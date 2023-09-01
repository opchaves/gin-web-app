package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/model"
)

type WorkspaceInput struct {
	name        string
	description string
	currency    string
	language    string
	userId      pgtype.UUID
}

type WorkspaceService interface {
	GetById(ctx context.Context, id string) (*model.Workspace, error)
	Create(ctx context.Context, data *WorkspaceInput) (*model.Workspace, error)
	BuildNewWorkspace(data *WorkspaceInput) (*model.CreateWorkspaceParams, error)
}

type workspaceService struct {
	Q      *model.Queries
	Logger *slog.Logger
	Db     *pgxpool.Pool
}

type WorkspaceServiceConfig struct {
	Q      *model.Queries
	Logger *slog.Logger
	Db     *pgxpool.Pool
}

func NewWorkspaceService(c *UserServiceConfig) WorkspaceService {
	return &workspaceService{
		Q:      c.Q,
		Logger: c.Logger,
		Db:     c.Db,
	}
}

// GetById implements WorkspaceService.
func (s *workspaceService) GetById(ctx context.Context, id string) (*model.Workspace, error) {
	panic("unimplemented")
}

// Create implements WorkspaceService.
func (s *workspaceService) Create(ctx context.Context, data *WorkspaceInput) (*model.Workspace, error) {
	// userId, err := utils.ConvertToUUID(data.userId)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// newWorkspace := model.CreateWorkspaceParams{
	// 	Name:        data.name,
	// 	Description: pgtype.Text{String: data.description, Valid: true},
	// 	Currency:    data.currency,
	// 	Language:    data.language,
	// 	UserID:      *userId,
	// }
	panic("unimplemented")
}

// BuildWorkspace implements WorkspaceService.
func (*workspaceService) BuildNewWorkspace(data *WorkspaceInput) (*model.CreateWorkspaceParams, error) {
	// userId, err := utils.ConvertToUUID(data.userId)
	// if err != nil {
	// 	return nil, err
	// }

	newWorkspace := model.CreateWorkspaceParams{
		Name:        data.name,
		Description: pgtype.Text{String: data.description, Valid: true},
		Currency:    data.currency,
		Language:    data.language,
		UserID:      data.userId,
	}

	return &newWorkspace, nil
}
