package preview

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/html"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

const (
	previewCachePrefix = "preview:"
	previewCacheTTL    = 24 * time.Hour
	maxBodyBytes       = 512 * 1024 // 512 KB
)

var previewHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return http.ErrUseLastResponse
		}
		return nil
	},
}

// ── Types ─────────────────────────────────────────────────────────────────────

type PreviewData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	SiteName    string `json:"site_name"`
	Favicon     string `json:"favicon"`
}

// ── Interface ─────────────────────────────────────────────────────────────────

type PreviewServiceI interface {
	GetPreview(ctx context.Context, linkID uuid.UUID, destURL string) (*PreviewData, *dto.ServiceError)
}

// ── Service ───────────────────────────────────────────────────────────────────

type previewService struct {
	cache valkeycompat.Cmdable
}

func NewPreviewService(cache valkeycompat.Cmdable) PreviewServiceI {
	return &previewService{cache: cache}
}

func (svc *previewService) GetPreview(ctx context.Context, linkID uuid.UUID, destURL string) (*PreviewData, *dto.ServiceError) {
	cacheKey := fmt.Sprintf("%s%s", previewCachePrefix, linkID.String())

	// Try Valkey cache first
	cached, err := svc.cache.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var data PreviewData
		if jsonErr := json.Unmarshal([]byte(cached), &data); jsonErr == nil {
			return &data, nil
		}
	}

	// Fetch OG data from URL
	data, fetchErr := fetchOGData(destURL)
	if fetchErr != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch link preview")
	}

	// Cache the result
	if b, jsonErr := json.Marshal(data); jsonErr == nil {
		svc.cache.Set(ctx, cacheKey, string(b), previewCacheTTL)
	}

	return data, nil
}

// ── OG fetcher & parser ───────────────────────────────────────────────────────

func fetchOGData(rawURL string) (*PreviewData, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return &PreviewData{}, nil
	}
	req.Header.Set("User-Agent", "Shortlink-Preview/1.0 (+https://shortlink.app)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,*/*;q=0.8")

	resp, err := previewHTTPClient.Do(req)
	if err != nil {
		return &PreviewData{}, nil
	}
	defer resp.Body.Close()

	body := io.LimitReader(resp.Body, maxBodyBytes)
	return parseOGTags(body, rawURL), nil
}

func parseOGTags(r io.Reader, baseURL string) *PreviewData {
	data := &PreviewData{}
	tokenizer := html.NewTokenizer(r)
	titleNext := false

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}

		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			tok := tokenizer.Token()
			name := strings.ToLower(tok.Data)

			switch name {
			case "body":
				return data // no need to parse beyond head

			case "title":
				titleNext = true

			case "meta":
				property := strings.ToLower(attrVal(tok, "property"))
				nameAttr := strings.ToLower(attrVal(tok, "name"))
				content := attrVal(tok, "content")

				switch property {
				case "og:title":
					if data.Title == "" {
						data.Title = content
					}
				case "og:description":
					if data.Description == "" {
						data.Description = content
					}
				case "og:image":
					if data.ImageURL == "" {
						data.ImageURL = resolveURL(baseURL, content)
					}
				case "og:site_name":
					if data.SiteName == "" {
						data.SiteName = content
					}
				}

				switch nameAttr {
				case "description":
					if data.Description == "" {
						data.Description = content
					}
				case "twitter:title":
					if data.Title == "" {
						data.Title = content
					}
				case "twitter:description":
					if data.Description == "" {
						data.Description = content
					}
				case "twitter:image":
					if data.ImageURL == "" {
						data.ImageURL = resolveURL(baseURL, content)
					}
				}

			case "link":
				rel := strings.ToLower(attrVal(tok, "rel"))
				href := attrVal(tok, "href")
				if (rel == "icon" || rel == "shortcut icon") && href != "" && data.Favicon == "" {
					data.Favicon = resolveURL(baseURL, href)
				}
			}

		case html.TextToken:
			if titleNext && data.Title == "" {
				data.Title = strings.TrimSpace(string(tokenizer.Text()))
				titleNext = false
			}

		case html.EndTagToken:
			tok := tokenizer.Token()
			if strings.ToLower(tok.Data) == "head" {
				return data
			}
		}
	}
	return data
}

func attrVal(tok html.Token, key string) string {
	for _, a := range tok.Attr {
		if strings.ToLower(a.Key) == key {
			return a.Val
		}
	}
	return ""
}

// resolveURL turns a potentially relative URL into an absolute one using baseURL.
func resolveURL(base, ref string) string {
	if ref == "" {
		return ""
	}
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		return ref
	}
	if strings.HasPrefix(ref, "//") {
		if strings.HasPrefix(base, "https") {
			return "https:" + ref
		}
		return "http:" + ref
	}

	// Determine the origin of base
	for _, scheme := range []string{"https://", "http://"} {
		if !strings.HasPrefix(base, scheme) {
			continue
		}
		rest := base[len(scheme):]
		slashIdx := strings.Index(rest, "/")
		var origin string
		if slashIdx == -1 {
			origin = scheme + rest
		} else {
			origin = scheme + rest[:slashIdx]
		}
		if strings.HasPrefix(ref, "/") {
			return origin + ref
		}
		// Relative to current path
		lastSlash := strings.LastIndex(base, "/")
		if lastSlash >= len(scheme) {
			return base[:lastSlash+1] + ref
		}
		return origin + "/" + ref
	}
	return ref
}
