package click

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/util"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

const clickQueueKey = "queue:clicks"

// ClickPayload is the data enqueued to Valkey for async processing.
type ClickPayload struct {
	LinkID     uuid.UUID `json:"link_id"`
	ClickedAt  time.Time `json:"clicked_at"`
	IPHash     string    `json:"ip_hash"`
	Country    string    `json:"country"`
	City       string    `json:"city"`
	DeviceType string    `json:"device_type"`
	OS         string    `json:"os"`
	Browser    string    `json:"browser"`
	Referrer   string    `json:"referrer"`
	Source     string    `json:"source"`
}

type ClickServiceI interface {
	EnqueueClick(ctx context.Context, linkID uuid.UUID, ip, userAgent, referrer, source string)
}

type clickService struct {
	cache valkeycompat.Cmdable
}

func NewClickService(cache valkeycompat.Cmdable) ClickServiceI {
	return &clickService{cache: cache}
}

func (s *clickService) EnqueueClick(ctx context.Context, linkID uuid.UUID, ip, userAgent, referrer, source string) {
	deviceInfo := util.ParseUserAgent(userAgent)
	countryCode, city := util.LookupGeo(ip)

	payload := ClickPayload{
		LinkID:     linkID,
		ClickedAt:  time.Now().UTC(),
		IPHash:     util.HashIP(ip),
		Country:    countryCode,
		City:       city,
		DeviceType: deviceInfo.DeviceType,
		OS:         deviceInfo.OS,
		Browser:    deviceInfo.Browser,
		Referrer:   referrer,
		Source:     source,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal click payload", zap.Error(err))
		return
	}

	if err := s.cache.RPush(ctx, clickQueueKey, string(data)).Err(); err != nil {
		logger.Error("Failed to enqueue click event", zap.Error(err))
	}
}

// ToClickEvent converts a ClickPayload to a model.ClickEvent for database insertion.
func ToClickEvent(p ClickPayload) *model.ClickEvent {
	return &model.ClickEvent{
		LinkID:     p.LinkID,
		ClickedAt:  p.ClickedAt,
		IPHash:     p.IPHash,
		Country:    p.Country,
		City:       p.City,
		DeviceType: p.DeviceType,
		OS:         p.OS,
		Browser:    p.Browser,
		Referrer:   p.Referrer,
		Source:     p.Source,
	}
}
