package click

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
	"github.com/shafikshaon/linkbee/util"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

const bioClickQueueKey = "queue:bio_clicks"

// BioClickPayload is the data enqueued to Valkey for async bio link click processing.
type BioClickPayload struct {
	BioLinkID  uuid.UUID `json:"bio_link_id"`
	BioPageID  uuid.UUID `json:"bio_page_id"`
	ClickedAt  time.Time `json:"clicked_at"`
	IPHash     string    `json:"ip_hash"`
	Country    string    `json:"country"`
	City       string    `json:"city"`
	DeviceType string    `json:"device_type"`
	OS         string    `json:"os"`
	Browser    string    `json:"browser"`
	Referrer   string    `json:"referrer"`
}

// BioClickServiceI enqueues click events for bio page links.
type BioClickServiceI interface {
	EnqueueBioClick(ctx context.Context, username string, bioLinkID uuid.UUID, ip, userAgent, referrer string)
}

type bioClickService struct {
	cache   valkeycompat.Cmdable
	bioRepo repository.BioRepositoryI
}

func NewBioClickService(cache valkeycompat.Cmdable, bioRepo repository.BioRepositoryI) BioClickServiceI {
	return &bioClickService{cache: cache, bioRepo: bioRepo}
}

// EnqueueBioClick resolves the bio page by username, validates the link belongs to
// that page, enriches the click with geo/UA data, and pushes it to the Valkey queue.
// Designed to be called in a goroutine — all errors are logged and swallowed.
func (s *bioClickService) EnqueueBioClick(ctx context.Context, username string, bioLinkID uuid.UUID, ip, userAgent, referrer string) {
	// Validate: page exists, is published, and owns the link.
	page, err := s.bioRepo.GetByUsername(ctx, username)
	if err != nil {
		return
	}
	var found bool
	for _, l := range page.Links {
		if l.ID == bioLinkID {
			found = true
			break
		}
	}
	if !found {
		return
	}

	deviceInfo := util.ParseUserAgent(userAgent)
	countryCode, city := util.LookupGeo(ip)

	payload := BioClickPayload{
		BioLinkID:  bioLinkID,
		BioPageID:  page.ID,
		ClickedAt:  time.Now().UTC(),
		IPHash:     util.HashIP(ip),
		Country:    countryCode,
		City:       city,
		DeviceType: deviceInfo.DeviceType,
		OS:         deviceInfo.OS,
		Browser:    deviceInfo.Browser,
		Referrer:   referrer,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal bio click payload", zap.Error(err))
		return
	}
	if err := s.cache.RPush(ctx, bioClickQueueKey, string(data)).Err(); err != nil {
		logger.Error("Failed to enqueue bio click event", zap.Error(err))
	}
}

// ToBioLinkClickEvent converts a BioClickPayload to a model.BioLinkClickEvent for DB insertion.
func ToBioLinkClickEvent(p BioClickPayload) *model.BioLinkClickEvent {
	return &model.BioLinkClickEvent{
		BioLinkID:  p.BioLinkID,
		BioPageID:  p.BioPageID,
		ClickedAt:  p.ClickedAt,
		IPHash:     p.IPHash,
		Country:    p.Country,
		City:       p.City,
		DeviceType: p.DeviceType,
		OS:         p.OS,
		Browser:    p.Browser,
		Referrer:   p.Referrer,
	}
}
