package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/model"
	"github.com/opchaves/gin-web-app/app/model/apperrors"
	"github.com/opchaves/gin-web-app/app/utils"
)

// type UserService interface {
// 	Get(ctx context.Context, id int32) (*User, error)
// 	GetByEmail(ctx context.Context, email string) (*User, error)
// 	Register(ctx context.Context, input *InsertUserParams) (*User, error)
// 	Login(ctx context.Context, email, password string) (*User, error)
// 	UpdateAccount(ctx context.Context, user *User) error
// 	IsEmailAlreadyInUse(ctx context.Context, email string) bool
// 	ChangePassword(ctx context.Context, currentPassword, newPassword string, user *User) error
// 	ForgotPassword(ctx context.Context, user *User) error
// 	ResetPassword(ctx context.Context, password string, token string) (*User, error)
// }

type RegisterInput struct {
	// Must be unique
	Email string `json:"email" binding:"required,email"`
	// Min 2, max 30 characters.
	FirstName string `json:"first_name" binding:"required,min=2,max=30"`
	// Min 2, max 30 characters.
	LastName string `json:"last_name" binding:"required,min=2,max=30"`
	// Min 10, max 100 characters.
	Password string `json:"password" binding:"required,min=10,max=100"`
} //@name RegisterRequest

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
} //@name LoginInput

type RegisterResponse struct {
	*model.User
	Password  bool `json:"password,omitempty"`
	LastLogin bool `json:"last_login,omitempty"`
	DeletedAt bool `json:"deleted_at,omitempty"`
} //@name RegisterResponse

type UserService interface {
	GetById(ctx context.Context, id string) (*RegisterResponse, error)
	Register(ctx context.Context, user *RegisterInput) (*RegisterResponse, error)
	Login(ctx context.Context, input *LoginInput) (*RegisterResponse, error)
}

type userService struct {
	Q      *model.Queries
	Logger *slog.Logger
	Db     *pgxpool.Pool
}

func NewUserService(c *ServiceConfig) UserService {
	return &userService{
		Q:      c.Q,
		Logger: c.Logger,
		Db:     c.Db,
	}
}

// GetById implements UserService.
func (us *userService) GetById(ctx context.Context, id string) (*RegisterResponse, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := us.Q.GetUserById(ctx, uuid)

	if err != nil {
		return nil, err
	}

	return &RegisterResponse{User: user}, err
}

// Register implements UserService.
func (us *userService) Register(ctx context.Context, data *RegisterInput) (*RegisterResponse, error) {
	hashedPassword, err := utils.HashPassword(data.Password)

	if err != nil {
		us.Logger.Error("unable to hash password", slog.Any("error", err))
		return nil, err
	}

	var lastLogin pgtype.Timestamp
	lastLogin.Scan(time.Now())

	newUser := model.CreateUserParams{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  hashedPassword,
		Active:    true,
		Role:      "user",
		LastLogin: lastLogin,
	}

	tx, err := us.Db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qTx := us.Q.WithTx(tx)
	user, err := qTx.CreateUser(ctx, newUser)

	if isDuplicateKeyError(err) {
		us.Logger.Warn("failed to register user", slog.Any("error", err))
		err = apperrors.NewBadRequest(apperrors.DuplicateEmail)
	}

	if err != nil {
		return nil, err
	}

	workspaceName := fmt.Sprintf("%s's workspace", user.FirstName)
	newWorkspace := model.CreateWorkspaceParams{
		Name:        workspaceName,
		Description: pgtype.Text{String: workspaceName, Valid: true},
		Currency:    "usd",
		Language:    "en-us",
		UserID:      user.ID,
	}

	_, err = qTx.CreateWorkspace(ctx, newWorkspace)
	if err != nil {
		us.Logger.Error("failed to create workspace", slog.String("userId", user.ID.String()))
		return nil, err
	}
	us.Logger.Info("User workspace created", slog.String("userId", user.ID.String()))

	tx.Commit(ctx)

	return &RegisterResponse{User: user}, err
	// TODO future: send email to verify account.
	// TODO when user verifies account, then create workspace, default accounts and categories
}

func (us *userService) Login(ctx context.Context, input *LoginInput) (*RegisterResponse, error) {
	user, err := us.Q.GetUserByEmail(ctx, input.Email)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	match, err := utils.ComparePasswords(user.Password, input.Password)

	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return &RegisterResponse{User: user}, err
}

// isDuplicateKeyError checks if the provided error is a PostgreSQL duplicate key error
func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
