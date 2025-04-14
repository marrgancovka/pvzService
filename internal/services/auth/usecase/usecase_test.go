package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/services/auth/mocks"
	"github.com/marrgancovka/pvzService/pkg/hasher"
	"go.uber.org/mock/gomock"
	"log/slog"
	"os"
	"testing"

	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_DummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockJWT := mocks.NewMockJWTer(ctrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	uc := Usecase{
		log:  log,
		repo: mockRepo,
		jwt:  mockJWT,
	}

	tests := []struct {
		name       string
		role       *models.DummyLogin
		setupMocks func()
		wantToken  string
		wantErr    error
	}{
		{
			name: "success",
			role: &models.DummyLogin{Role: models.RoleEmployee},
			setupMocks: func() {
				mockJWT.EXPECT().
					GenerateJWT(&models.TokenPayload{Role: models.RoleEmployee}).
					Return(&models.Token{Token: "testToken"}, nil)
			},
			wantToken: "testToken",
			wantErr:   nil,
		},
		{
			name:       "invalid role",
			role:       &models.DummyLogin{Role: "invalid_role"},
			setupMocks: func() {},
			wantToken:  "",
			wantErr:    auth.ErrIncorrectRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			token, err := uc.DummyLogin(context.Background(), tt.role)

			assert.Equal(t, tt.wantToken, token)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if !errors.Is(tt.wantErr, auth.ErrIncorrectRole) {
					assert.EqualError(t, err, tt.wantErr.Error())
				} else {
					assert.True(t, errors.Is(err, auth.ErrIncorrectRole))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockJWT := mocks.NewMockJWTer(ctrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	uc := Usecase{
		log:  log,
		repo: mockRepo,
		jwt:  mockJWT,
	}

	tests := []struct {
		name       string
		userData   *models.Users
		setupMocks func()
		wantToken  string
		wantErr    error
	}{
		{
			name: "successful login",
			userData: &models.Users{
				Email:    "test@example.com",
				Password: "password",
			},
			setupMocks: func() {
				mockRepo.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&models.Users{
						Email:    "test@example.com",
						Password: hasher.GenerateHashString("password"),
						ID:       uuid.New(),
						Role:     "employee",
					}, nil)
				mockJWT.EXPECT().
					GenerateJWT(gomock.Any()).
					Return(&models.Token{Token: "testToken"}, nil)
			},
			wantToken: "testToken",
			wantErr:   nil,
		},
		{
			name: "user not found",
			userData: &models.Users{
				Email:    "nonexist@example.com",
				Password: "password",
			},
			setupMocks: func() {
				mockRepo.EXPECT().
					GetUserByEmail(gomock.Any(), "nonexist@example.com").
					Return(nil, auth.ErrUserNotFound)
			},
			wantToken: "",
			wantErr:   auth.ErrUserNotFound,
		},
		{
			name: "incorrect password",
			userData: &models.Users{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			setupMocks: func() {
				mockRepo.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&models.Users{
						Email:    "test@example.com",
						Password: "hashed_password",
					}, nil)
			},
			wantToken: "",
			wantErr:   auth.ErrIncorrectPasswordOrEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			token, err := uc.Login(context.Background(), tt.userData)
			fmt.Println(err)
			assert.Equal(t, tt.wantToken, token)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if !errors.Is(tt.wantErr, auth.ErrIncorrectPasswordOrEmail) {
					assert.EqualError(t, err, tt.wantErr.Error())
				} else {
					assert.True(t, errors.Is(err, auth.ErrIncorrectPasswordOrEmail))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockJWT := mocks.NewMockJWTer(ctrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	uc := Usecase{
		log:  log,
		repo: mockRepo,
		jwt:  mockJWT,
	}

	userID := uuid.New()

	tests := []struct {
		name       string
		userData   *models.Users
		setupMocks func()
		wantToken  string
		wantErr    error
	}{
		{
			name: "successful",
			userData: &models.Users{
				ID:       userID,
				Email:    "example@example.com",
				Password: "password",
				Role:     "employee",
			},
			setupMocks: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).Return(&models.Users{
					ID:       userID,
					Email:    "example@example.com",
					Password: hasher.GenerateHashString("password"),
					Role:     "employee",
				}, nil)
				mockJWT.EXPECT().
					GenerateJWT(&models.TokenPayload{
						ID:   userID,
						Role: "employee",
					}).
					Return(&models.Token{Token: "testToken"}, nil)
			},
			wantToken: "testToken",
			wantErr:   nil,
		},
		{
			name: "invalid role",
			userData: &models.Users{
				Email:    "example@example.com",
				Password: "password",
				Role:     "unknown",
			},
			setupMocks: func() {},
			wantToken:  "",
			wantErr:    auth.ErrIncorrectRole,
		},
		{
			name: "user already exists",
			userData: &models.Users{
				Email:    "exist@example.com",
				Password: "password",
				Role:     "employee",
			},
			setupMocks: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, auth.ErrUserAlreadyExists)
			},
			wantToken: "",
			wantErr:   auth.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			token, err := uc.Register(context.Background(), tt.userData)

			assert.Equal(t, tt.wantToken, token)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if !errors.Is(tt.wantErr, auth.ErrIncorrectRole) {
					assert.EqualError(t, err, tt.wantErr.Error())
				} else {
					assert.True(t, errors.Is(err, auth.ErrIncorrectRole))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
