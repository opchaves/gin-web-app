package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

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
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=10,max=50"`
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserService interface {
	GetById(ctx context.Context, id string) (*model.User, error)
	Register(ctx context.Context, user *RegisterInput) (*model.User, error)
	GetRegisterResponse(user *model.User) *RegisterResponse
}

type userService struct {
	Q      *model.Queries
	Logger *slog.Logger
	Db     *pgxpool.Pool
}

type UserServiceConfig struct {
	Q      *model.Queries
	Logger *slog.Logger
	Db     *pgxpool.Pool
}

func NewUserService(c *UserServiceConfig) UserService {
	return &userService{
		Q:      c.Q,
		Logger: c.Logger,
		Db:     c.Db,
	}
}

// GetById implements UserService.
func (us *userService) GetById(ctx context.Context, id string) (*model.User, error) {
	uuid, err := utils.ConvertToUUID(id)
	if err != nil {
		return nil, err
	}

	return us.Q.GetUserById(ctx, *uuid)
}

// Register implements UserService.
func (us *userService) Register(ctx context.Context, data *RegisterInput) (*model.User, error) {
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
		us.Logger.Error("failed to create workspace", slog.String("userId", utils.UUIDtoString(user.ID)))
		return nil, err
	}
	us.Logger.Info("User workspace created", slog.String("userId", utils.UUIDtoString(user.ID)))

	tx.Commit(ctx)

	return user, err
	// TODO future: send email to verify account.
	// TODO when user verifies account, then create workspace, default accounts and categories
}

// GetRegisterResponse implements UserService.
func (s *userService) GetRegisterResponse(user *model.User) *RegisterResponse {
	var r RegisterResponse

	r.ID = utils.UUIDtoString(user.ID)
	r.FirstName = user.FirstName
	r.LastName = user.LastName
	r.Email = user.Email
	r.Role = user.Role
	r.CreatedAt = user.CreatedAt.Time
	r.UpdatedAt = user.UpdatedAt.Time

	return &r
}

// isDuplicateKeyError checks if the provided error is a PostgreSQL duplicate key error
func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
