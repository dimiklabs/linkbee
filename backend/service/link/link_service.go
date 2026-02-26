package link

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/request"
	"github.com/shafikshaon/shortlink/response"
	"github.com/shafikshaon/shortlink/util"
)

type LinkServiceI interface {
	CreateLink(ctx context.Context, userID uuid.UUID, req *request.CreateLinkRequest) (*response.LinkResponse, *dto.ServiceError)
	GetLink(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*response.LinkResponse, *dto.ServiceError)
	ListLinks(ctx context.Context, userID uuid.UUID, page, limit int, search string, folderID *uuid.UUID, starred *bool, healthStatus string) (*response.LinkListResponse, *dto.ServiceError)
	UpdateLink(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *request.UpdateLinkRequest) (*response.LinkResponse, *dto.ServiceError)
	DeleteLink(ctx context.Context, id uuid.UUID, userID uuid.UUID) *dto.ServiceError
	ToggleStar(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*response.LinkResponse, *dto.ServiceError)
	ImportLinks(ctx context.Context, userID uuid.UUID, r io.Reader) (*response.ImportLinksResponse, *dto.ServiceError)
	CheckLinkHealth(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*response.LinkResponse, *dto.ServiceError)
}

// healthHTTPClient is shared across all on-demand health checks.
var healthHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return http.ErrUseLastResponse
		}
		return nil
	},
}

type linkService struct {
	linkRepo repository.LinkRepositoryI
	appCfg   *config.AppConfig
	linkCfg  *config.LinkConfig
}

func NewLinkService(linkRepo repository.LinkRepositoryI, appCfg *config.AppConfig, linkCfg *config.LinkConfig) LinkServiceI {
	return &linkService{
		linkRepo: linkRepo,
		appCfg:   appCfg,
		linkCfg:  linkCfg,
	}
}

func (s *linkService) CreateLink(ctx context.Context, userID uuid.UUID, req *request.CreateLinkRequest) (*response.LinkResponse, *dto.ServiceError) {
	// Validate destination URL
	if _, err := url.ParseRequestURI(req.DestinationURL); err != nil {
		return nil, dto.NewBadRequestError(constant.ErrCodeInvalidURL, constant.ErrMsgInvalidURL)
	}

	// Determine slug
	slug := req.Slug
	if slug == "" {
		// Generate a unique slug
		for attempt := 0; attempt < 5; attempt++ {
			generated, err := util.GenerateSlug(s.linkCfg.SlugLength)
			if err != nil {
				logger.ErrorCtx(ctx, "Failed to generate slug", zap.Error(err))
				return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
			}
			exists, err := s.linkRepo.SlugExists(ctx, generated)
			if err != nil {
				return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
			}
			if !exists {
				slug = generated
				break
			}
		}
		if slug == "" {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "Failed to generate unique slug")
		}
	} else {
		// Check custom slug availability
		exists, err := s.linkRepo.SlugExists(ctx, slug)
		if err != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		if exists {
			return nil, dto.NewConflictError(constant.ErrCodeSlugTaken, constant.ErrMsgSlugTaken)
		}
	}

	link := &model.Link{
		UserID:         userID,
		Slug:           slug,
		DestinationURL: req.DestinationURL,
		Title:          req.Title,
		RedirectType:   int16(s.linkCfg.DefaultRedirectType),
		IsActive:       true,
		UTMSource:      req.UTMSource,
		UTMMedium:      req.UTMMedium,
		UTMCampaign:    req.UTMCampaign,
	}

	if req.FolderID != nil && *req.FolderID != "" {
		folderID, err := uuid.Parse(*req.FolderID)
		if err == nil {
			link.FolderID = &folderID
		}
	}

	if req.RedirectType != nil {
		link.RedirectType = *req.RedirectType
	}

	if req.MaxClicks != nil {
		link.MaxClicks = req.MaxClicks
	}

	if len(req.Tags) > 0 {
		link.Tags = req.Tags
	}

	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return nil, dto.NewBadRequestError(constant.ErrCodeValidationError, "expires_at must be in RFC3339 format")
		}
		link.ExpiresAt = &t
	}

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		link.PasswordHash = string(hash)
	}

	if err := s.linkRepo.Create(ctx, link); err != nil {
		logger.ErrorCtx(ctx, "Failed to create link", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return s.toLinkResponse(link), nil
}

func (s *linkService) GetLink(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*response.LinkResponse, *dto.ServiceError) {
	link, err := s.linkRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	if link.UserID != userID {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}

	return s.toLinkResponse(link), nil
}

func (s *linkService) ListLinks(ctx context.Context, userID uuid.UUID, page, limit int, search string, folderID *uuid.UUID, starred *bool, healthStatus string) (*response.LinkListResponse, *dto.ServiceError) {
	links, total, err := s.linkRepo.GetByUserID(ctx, userID, page, limit, search, folderID, starred, healthStatus)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	linkResponses := make([]*response.LinkResponse, len(links))
	for i := range links {
		linkResponses[i] = s.toLinkResponse(&links[i])
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &response.LinkListResponse{
		Links:      linkResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *linkService) UpdateLink(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *request.UpdateLinkRequest) (*response.LinkResponse, *dto.ServiceError) {
	link, err := s.linkRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	if link.UserID != userID {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}

	if req.DestinationURL != "" {
		link.DestinationURL = req.DestinationURL
	}
	if req.Title != "" {
		link.Title = req.Title
	}
	if req.RedirectType != nil {
		link.RedirectType = *req.RedirectType
	}
	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	}
	if req.MaxClicks != nil {
		link.MaxClicks = req.MaxClicks
	}
	if req.Tags != nil {
		link.Tags = req.Tags
	}
	if req.UTMSource != "" {
		link.UTMSource = req.UTMSource
	}
	if req.UTMMedium != "" {
		link.UTMMedium = req.UTMMedium
	}
	if req.UTMCampaign != "" {
		link.UTMCampaign = req.UTMCampaign
	}
	if req.ExpiresAt != nil {
		if *req.ExpiresAt == "" {
			link.ExpiresAt = nil
		} else {
			t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
			if err != nil {
				return nil, dto.NewBadRequestError(constant.ErrCodeValidationError, "expires_at must be in RFC3339 format")
			}
			link.ExpiresAt = &t
		}
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		link.PasswordHash = string(hash)
	}
	if req.FolderID != nil {
		if *req.FolderID == "" {
			link.FolderID = nil
		} else {
			folderID, err := uuid.Parse(*req.FolderID)
			if err == nil {
				link.FolderID = &folderID
			}
		}
	}

	if err := s.linkRepo.Update(ctx, link); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return s.toLinkResponse(link), nil
}

func (s *linkService) DeleteLink(ctx context.Context, id uuid.UUID, userID uuid.UUID) *dto.ServiceError {
	if err := s.linkRepo.Delete(ctx, id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

func (s *linkService) ToggleStar(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*response.LinkResponse, *dto.ServiceError) {
	isStarred, err := s.linkRepo.ToggleStar(ctx, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	link, err := s.linkRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	link.IsStarred = isStarred
	return s.toLinkResponse(link), nil
}

func (s *linkService) CheckLinkHealth(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*response.LinkResponse, *dto.ServiceError) {
	link, err := s.linkRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if link.UserID != userID {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}

	status, statusCode := probeURL(ctx, link.DestinationURL)

	if err := s.linkRepo.UpdateHealthStatus(ctx, id, status, statusCode); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	link.HealthStatus = status
	link.HealthStatusCode = statusCode
	now := time.Now()
	link.HealthCheckedAt = &now
	return s.toLinkResponse(link), nil
}

// probeURL performs a HEAD (then GET fallback) request and returns a health status + HTTP code.
func probeURL(ctx context.Context, rawURL string) (status string, statusCode int) {
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(checkCtx, http.MethodHead, rawURL, nil)
	if err != nil {
		return response.HealthStatusError, 0
	}
	req.Header.Set("User-Agent", "shortlink-health-monitor/1.0")

	resp, err := healthHTTPClient.Do(req)
	if err != nil {
		if checkCtx.Err() != nil {
			return response.HealthStatusTimeout, 0
		}
		return response.HealthStatusError, 0
	}
	resp.Body.Close()

	// Some servers reject HEAD — retry with GET
	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotImplemented {
		getReq, gErr := http.NewRequestWithContext(checkCtx, http.MethodGet, rawURL, nil)
		if gErr == nil {
			getReq.Header.Set("User-Agent", "shortlink-health-monitor/1.0")
			if getResp, gErr2 := healthHTTPClient.Do(getReq); gErr2 == nil {
				resp.StatusCode = getResp.StatusCode
				getResp.Body.Close()
			}
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return response.HealthStatusHealthy, resp.StatusCode
	}
	return response.HealthStatusUnhealthy, resp.StatusCode
}

const importMaxRows = 500

func (s *linkService) ImportLinks(ctx context.Context, userID uuid.UUID, r io.Reader) (*response.ImportLinksResponse, *dto.ServiceError) {
	csvReader := csv.NewReader(r)
	csvReader.TrimLeadingSpace = true

	header, err := csvReader.Read()
	if err != nil {
		return nil, dto.NewBadRequestError(constant.ErrCodeBadRequest, "Invalid CSV: could not read header row")
	}

	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[strings.ToLower(strings.TrimSpace(col))] = i
	}

	if _, ok := colIdx["destination_url"]; !ok {
		return nil, dto.NewBadRequestError(constant.ErrCodeBadRequest, "CSV must have a 'destination_url' column")
	}

	result := &response.ImportLinksResponse{
		Errors: []response.ImportLinkError{},
	}
	rowNum := 0

	for {
		record, readErr := csvReader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			rowNum++
			result.Total++
			result.Failed++
			result.Errors = append(result.Errors, response.ImportLinkError{
				Row:   rowNum,
				Error: "Parse error: " + readErr.Error(),
			})
			continue
		}

		// Skip blank rows
		allEmpty := true
		for _, cell := range record {
			if strings.TrimSpace(cell) != "" {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			continue
		}

		rowNum++
		if rowNum > importMaxRows {
			break
		}
		result.Total++

		destURL := csvField(record, colIdx, "destination_url")
		if destURL == "" {
			result.Failed++
			result.Errors = append(result.Errors, response.ImportLinkError{
				Row:   rowNum,
				Error: "destination_url is required",
			})
			continue
		}

		req := &request.CreateLinkRequest{
			DestinationURL: destURL,
			Slug:           csvField(record, colIdx, "slug"),
			Title:          csvField(record, colIdx, "title"),
		}

		if tagsStr := csvField(record, colIdx, "tags"); tagsStr != "" {
			for _, t := range strings.Split(tagsStr, ";") {
				if t = strings.TrimSpace(t); t != "" {
					req.Tags = append(req.Tags, t)
				}
			}
		}

		switch csvField(record, colIdx, "redirect_type") {
		case "301":
			rt := int16(301)
			req.RedirectType = &rt
		case "302":
			rt := int16(302)
			req.RedirectType = &rt
		}

		if fid := csvField(record, colIdx, "folder_id"); fid != "" {
			req.FolderID = &fid
		}

		if _, svcErr := s.CreateLink(ctx, userID, req); svcErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, response.ImportLinkError{
				Row:   rowNum,
				URL:   destURL,
				Error: svcErr.Description,
			})
		} else {
			result.Created++
		}
	}

	logger.InfoCtx(ctx, "CSV import complete",
		zap.String("user_id", userID.String()),
		zap.Int("total", result.Total),
		zap.Int("created", result.Created),
		zap.Int("failed", result.Failed))

	return result, nil
}

func csvField(record []string, colIdx map[string]int, col string) string {
	idx, ok := colIdx[col]
	if !ok || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}

func (s *linkService) toLinkResponse(link *model.Link) *response.LinkResponse {
	shortURL := fmt.Sprintf("%s/%s", s.appCfg.BaseDomain, link.Slug)

	// Append UTM params if set
	if link.UTMSource != "" || link.UTMMedium != "" || link.UTMCampaign != "" {
		parsed, err := url.Parse(link.DestinationURL)
		if err == nil {
			q := parsed.Query()
			if link.UTMSource != "" {
				q.Set("utm_source", link.UTMSource)
			}
			if link.UTMMedium != "" {
				q.Set("utm_medium", link.UTMMedium)
			}
			if link.UTMCampaign != "" {
				q.Set("utm_campaign", link.UTMCampaign)
			}
			parsed.RawQuery = q.Encode()
		}
	}

	tags := []string{}
	if link.Tags != nil {
		tags = link.Tags
	}

	return &response.LinkResponse{
		ID:               link.ID,
		FolderID:         link.FolderID,
		Slug:             link.Slug,
		ShortURL:         shortURL,
		DestinationURL:   link.DestinationURL,
		Title:            link.Title,
		ClickCount:       link.ClickCount,
		RedirectType:     link.RedirectType,
		IsActive:         link.IsActive,
		IsStarred:        link.IsStarred,
		IsSplitTest:      link.IsSplitTest,
		IsGeoRouting:     link.IsGeoRouting,
		HealthStatus:     link.HealthStatus,
		HealthStatusCode: link.HealthStatusCode,
		HealthCheckedAt:  link.HealthCheckedAt,
		Tags:             tags,
		HasPassword:      link.PasswordHash != "",
		ExpiresAt:        link.ExpiresAt,
		MaxClicks:        link.MaxClicks,
		UTMSource:        link.UTMSource,
		UTMMedium:        link.UTMMedium,
		UTMCampaign:      link.UTMCampaign,
		CreatedAt:        link.CreatedAt,
		UpdatedAt:        link.UpdatedAt,
	}
}
