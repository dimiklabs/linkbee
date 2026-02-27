package folder

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/request"
	"github.com/shafikshaon/shortlink/response"
)

type FolderServiceI interface {
	CreateFolder(ctx context.Context, userID uuid.UUID, req *request.CreateFolderRequest) (*response.FolderResponse, *dto.ServiceError)
	ListFolders(ctx context.Context, userID uuid.UUID) ([]response.FolderResponse, *dto.ServiceError)
	UpdateFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *request.UpdateFolderRequest) (*response.FolderResponse, *dto.ServiceError)
	DeleteFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID) *dto.ServiceError
}

type folderService struct {
	folderRepo     repository.FolderRepositoryI
	clickEventRepo repository.ClickEventRepositoryI
}

func NewFolderService(folderRepo repository.FolderRepositoryI, clickEventRepo repository.ClickEventRepositoryI) FolderServiceI {
	return &folderService{folderRepo: folderRepo, clickEventRepo: clickEventRepo}
}

func (s *folderService) CreateFolder(ctx context.Context, userID uuid.UUID, req *request.CreateFolderRequest) (*response.FolderResponse, *dto.ServiceError) {
	color := req.Color
	if color == "" {
		color = "#635bff"
	}
	folder := &model.Folder{
		UserID: userID,
		Name:   req.Name,
		Color:  color,
	}
	if err := s.folderRepo.Create(ctx, folder); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return toFolderResponse(folder), nil
}

func (s *folderService) ListFolders(ctx context.Context, userID uuid.UUID) ([]response.FolderResponse, *dto.ServiceError) {
	folders, err := s.folderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	result := make([]response.FolderResponse, len(folders))
	for i := range folders {
		r := toFolderResponse(&folders[i])
		if count, cErr := s.clickEventRepo.GetClickCountByFolderID(ctx, folders[i].ID); cErr == nil {
			r.ClickCount = count
		}
		result[i] = *r
	}
	return result, nil
}

func (s *folderService) UpdateFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *request.UpdateFolderRequest) (*response.FolderResponse, *dto.ServiceError) {
	folder, err := s.folderRepo.GetByID(ctx, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError("FOLDER_NOT_FOUND", "Folder not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if req.Name != "" {
		folder.Name = req.Name
	}
	if req.Color != "" {
		folder.Color = req.Color
	}
	if err := s.folderRepo.Update(ctx, folder); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return toFolderResponse(folder), nil
}

func (s *folderService) DeleteFolder(ctx context.Context, id uuid.UUID, userID uuid.UUID) *dto.ServiceError {
	if err := s.folderRepo.Delete(ctx, id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError("FOLDER_NOT_FOUND", "Folder not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

func toFolderResponse(f *model.Folder) *response.FolderResponse {
	return &response.FolderResponse{
		ID:        f.ID,
		Name:      f.Name,
		Color:     f.Color,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
