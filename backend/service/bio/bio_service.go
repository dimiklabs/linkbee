package bio

import (
	"context"
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
)

// ── Response types ────────────────────────────────────────────────────────────

type BioLinkResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	IsActive bool   `json:"is_active"`
	Position int    `json:"position"`
}

type BioPageResponse struct {
	ID          string            `json:"id"`
	Username    string            `json:"username"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	AvatarURL   string            `json:"avatar_url"`
	Theme       string            `json:"theme"`
	IsPublished bool              `json:"is_published"`
	Links       []BioLinkResponse `json:"links"`
	CreatedAt   string            `json:"created_at"`
}

// ── Interface ─────────────────────────────────────────────────────────────────

type BioServiceI interface {
	GetOrCreate(ctx context.Context, userID uuid.UUID) (*BioPageResponse, *dto.ServiceError)
	Update(ctx context.Context, userID uuid.UUID, req UpdateBioRequest) (*BioPageResponse, *dto.ServiceError)
	GetPublic(ctx context.Context, username string) (*BioPageResponse, *dto.ServiceError)

	ListLinks(ctx context.Context, userID uuid.UUID) ([]BioLinkResponse, *dto.ServiceError)
	CreateLink(ctx context.Context, userID uuid.UUID, title, url string) (*BioLinkResponse, *dto.ServiceError)
	UpdateLink(ctx context.Context, userID uuid.UUID, linkID uuid.UUID, title, url string, isActive bool) (*BioLinkResponse, *dto.ServiceError)
	DeleteLink(ctx context.Context, userID uuid.UUID, linkID uuid.UUID) *dto.ServiceError
	ReorderLinks(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) *dto.ServiceError
}

type UpdateBioRequest struct {
	Username    string
	Title       string
	Description string
	AvatarURL   string
	Theme       string
	IsPublished bool
}

// ── Service ───────────────────────────────────────────────────────────────────

type BioService struct {
	repo repository.BioRepositoryI
}

func NewBioService(repo repository.BioRepositoryI) BioServiceI {
	return &BioService{repo: repo}
}

var usernameRe = regexp.MustCompile(`[^a-z0-9_]`)

func sanitizeUsername(s string) string {
	s = strings.ToLower(s)
	s = usernameRe.ReplaceAllString(s, "")
	if len(s) > 30 {
		s = s[:30]
	}
	return s
}

func (svc *BioService) generateUsername(ctx context.Context, base string) (string, error) {
	base = sanitizeUsername(base)
	if len(base) < 3 {
		base = "user"
	}

	exists, err := svc.repo.UsernameExists(ctx, base)
	if err != nil {
		return "", err
	}
	if !exists {
		return base, nil
	}

	// Append random 3-digit suffix
	buf := make([]byte, 2)
	for i := 0; i < 10; i++ {
		_, _ = rand.Read(buf)
		candidate := fmt.Sprintf("%s%d", base, int(buf[0])%900+100)
		exists, err = svc.repo.UsernameExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not generate unique username")
}

func toBioPageResponse(p *model.BioPage) *BioPageResponse {
	links := make([]BioLinkResponse, 0, len(p.Links))
	for _, l := range p.Links {
		links = append(links, BioLinkResponse{
			ID:       l.ID.String(),
			Title:    l.Title,
			URL:      l.URL,
			IsActive: l.IsActive,
			Position: l.Position,
		})
	}
	return &BioPageResponse{
		ID:          p.ID.String(),
		Username:    p.Username,
		Title:       p.Title,
		Description: p.Description,
		AvatarURL:   p.AvatarURL,
		Theme:       p.Theme,
		IsPublished: p.IsPublished,
		Links:       links,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toBioLinkResponse(l *model.BioLink) *BioLinkResponse {
	return &BioLinkResponse{
		ID:       l.ID.String(),
		Title:    l.Title,
		URL:      l.URL,
		IsActive: l.IsActive,
		Position: l.Position,
	}
}

func (svc *BioService) GetOrCreate(ctx context.Context, userID uuid.UUID) (*BioPageResponse, *dto.ServiceError) {
	page, err := svc.repo.GetByUserID(ctx, userID)
	if err == nil {
		return toBioPageResponse(page), nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}

	// Auto-generate username from first segment of UUID
	base := strings.ReplaceAll(userID.String()[:8], "-", "")
	username, genErr := svc.generateUsername(ctx, "user"+base)
	if genErr != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to generate username")
	}

	page = &model.BioPage{
		UserID:      userID,
		Username:    username,
		Title:       "My Links",
		Theme:       "light",
		IsPublished: false,
		Links:       []model.BioLink{},
	}
	if err := svc.repo.Create(ctx, page); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to create bio page")
	}
	return toBioPageResponse(page), nil
}

func (svc *BioService) Update(ctx context.Context, userID uuid.UUID, req UpdateBioRequest) (*BioPageResponse, *dto.ServiceError) {
	page, err := svc.repo.GetByUserID(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "bio page not found")
	}
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}

	if req.Username != "" && req.Username != page.Username {
		clean := sanitizeUsername(req.Username)
		if len(clean) < 3 {
			return nil, dto.NewBadRequestError(constant.ErrCodeBadRequest, "username must be at least 3 characters (letters, numbers, underscores only)")
		}
		exists, checkErr := svc.repo.UsernameExists(ctx, clean)
		if checkErr != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to check username")
		}
		if exists {
			return nil, dto.NewConflictError(constant.ErrCodeConflict, "username already taken")
		}
		page.Username = clean
	}

	if req.Title != "" {
		page.Title = req.Title
	}
	page.Description = req.Description
	page.AvatarURL = req.AvatarURL
	if req.Theme != "" {
		page.Theme = req.Theme
	}
	page.IsPublished = req.IsPublished

	if err := svc.repo.Update(ctx, page); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to update bio page")
	}

	// Reload with links
	page, err = svc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to reload bio page")
	}
	return toBioPageResponse(page), nil
}

func (svc *BioService) GetPublic(ctx context.Context, username string) (*BioPageResponse, *dto.ServiceError) {
	page, err := svc.repo.GetByUsername(ctx, username)
	if err == gorm.ErrRecordNotFound {
		return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "bio page not found")
	}
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}
	return toBioPageResponse(page), nil
}

func (svc *BioService) ListLinks(ctx context.Context, userID uuid.UUID) ([]BioLinkResponse, *dto.ServiceError) {
	page, err := svc.repo.GetByUserID(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "bio page not found")
	}
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}

	links, err := svc.repo.ListLinks(ctx, page.ID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to list links")
	}
	result := make([]BioLinkResponse, 0, len(links))
	for i := range links {
		result = append(result, *toBioLinkResponse(&links[i]))
	}
	return result, nil
}

func (svc *BioService) CreateLink(ctx context.Context, userID uuid.UUID, title, url string) (*BioLinkResponse, *dto.ServiceError) {
	page, err := svc.repo.GetByUserID(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "bio page not found")
	}
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}

	links, err := svc.repo.ListLinks(ctx, page.ID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to count links")
	}

	link := &model.BioLink{
		BioPageID: page.ID,
		Title:     title,
		URL:       url,
		IsActive:  true,
		Position:  len(links),
	}
	if err := svc.repo.CreateLink(ctx, link); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to create link")
	}
	return toBioLinkResponse(link), nil
}

func (svc *BioService) UpdateLink(ctx context.Context, userID uuid.UUID, linkID uuid.UUID, title, url string, isActive bool) (*BioLinkResponse, *dto.ServiceError) {
	page, err := svc.repo.GetByUserID(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "bio page not found")
	}
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}

	link, err := svc.repo.GetLink(ctx, linkID, page.ID)
	if err == gorm.ErrRecordNotFound {
		return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "link not found")
	}
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch link")
	}

	if title != "" {
		link.Title = title
	}
	if url != "" {
		link.URL = url
	}
	link.IsActive = isActive

	if err := svc.repo.UpdateLink(ctx, link); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to update link")
	}
	return toBioLinkResponse(link), nil
}

func (svc *BioService) DeleteLink(ctx context.Context, userID uuid.UUID, linkID uuid.UUID) *dto.ServiceError {
	page, err := svc.repo.GetByUserID(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		return dto.NewNotFoundError(constant.ErrCodeNotFound, "bio page not found")
	}
	if err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}

	if err := svc.repo.DeleteLink(ctx, linkID, page.ID); err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to delete link")
	}
	return nil
}

func (svc *BioService) ReorderLinks(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) *dto.ServiceError {
	page, err := svc.repo.GetByUserID(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		return dto.NewNotFoundError(constant.ErrCodeNotFound, "bio page not found")
	}
	if err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch bio page")
	}

	if err := svc.repo.ReorderLinks(ctx, page.ID, ids); err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to reorder links")
	}
	return nil
}
