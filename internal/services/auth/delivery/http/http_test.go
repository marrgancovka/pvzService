package http

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	authMocks "github.com/marrgancovka/pvzService/internal/services/auth/mocks"
	"github.com/marrgancovka/pvzService/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_DummyLogin(t *testing.T) {
	type mockBehavior func(m *authMocks.MockUsecase, role *models.DummyLogin)
	testTable := []struct {
		name         string
		inputBody    string
		inputRole    *models.DummyLogin
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name:      "success moderator",
			inputBody: `{"role":"moderator"}`,
			inputRole: &models.DummyLogin{
				Role: "moderator",
			},
			mockBehavior: func(m *authMocks.MockUsecase, role *models.DummyLogin) {
				m.EXPECT().DummyLogin(gomock.Any(), role).Return("testToken", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"testToken"`,
		},
		{
			name:      "success employee",
			inputBody: `{"role":"employee"}`,
			inputRole: &models.DummyLogin{
				Role: "employee",
			},
			mockBehavior: func(m *authMocks.MockUsecase, role *models.DummyLogin) {
				m.EXPECT().DummyLogin(gomock.Any(), role).Return("testToken", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"testToken"`,
		},
		{
			name:      "empty body",
			inputBody: `{}`,
			inputRole: &models.DummyLogin{
				Role: "",
			},
			mockBehavior: func(m *authMocks.MockUsecase, role *models.DummyLogin) {
				m.EXPECT().DummyLogin(gomock.Any(), role).Return("", auth.ErrIncorrectRole)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrIncorrectRole.Error() + `"}`,
		},
		{
			name:         "invalid JSON",
			inputBody:    `{"role":123}`,
			inputRole:    nil,
			mockBehavior: func(m *authMocks.MockUsecase, role *models.DummyLogin) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg":"bad request"}`,
		},
		{
			name:      "incorrect role",
			inputBody: `{"role":"unknown"}`,
			inputRole: &models.DummyLogin{
				Role: "unknown",
			},
			mockBehavior: func(m *authMocks.MockUsecase, role *models.DummyLogin) {
				m.EXPECT().DummyLogin(gomock.Any(), role).Return("", auth.ErrIncorrectRole)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrIncorrectRole.Error() + `"}`,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecaseEmployee := authMocks.NewMockUsecase(ctrl)
			tt.mockBehavior(mockUsecaseEmployee, tt.inputRole)

			handler := &Handler{
				usecase: mockUsecaseEmployee,
				logger:  logger.SetupLogger(),
			}

			router := mux.NewRouter()
			router.HandleFunc("/dummyLogin", handler.DummyLogin).Methods(http.MethodPost)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(tt.inputBody))

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

		})
	}
}

func TestHandler_Login(t *testing.T) {
	type mockBehavior func(m *authMocks.MockUsecase, user *models.Users)

	testTable := []struct {
		name         string
		inputBody    string
		inputUser    *models.Users
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name:      "success login",
			inputBody: `{"email":"test@example.com", "password":"password123"}`,
			inputUser: &models.Users{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {
				m.EXPECT().Login(gomock.Any(), user).Return("valid_token", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"valid_token"`,
		},
		{
			name:      "empty body",
			inputBody: `{}`,
			inputUser: &models.Users{
				Email:    "",
				Password: "",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrBadRequest.Error() + `"}`,
		},
		{
			name:         "invalid JSON",
			inputBody:    `{"email":123}`,
			inputUser:    nil,
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrBadRequest.Error() + `"}`,
		},
		{
			name:      "user not found",
			inputBody: `{"email":"nonexistent@example.com", "password":"wrong"}`,
			inputUser: &models.Users{
				Email:    "nonexistent@example.com",
				Password: "wrong",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {
				m.EXPECT().Login(gomock.Any(), user).Return("", auth.ErrUserNotFound)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrIncorrectData.Error() + `"}`,
		},
		{
			name:      "incorrect password",
			inputBody: `{"email":"test@example.com", "password":"wrong"}`,
			inputUser: &models.Users{
				Email:    "test@example.com",
				Password: "wrong",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {
				m.EXPECT().Login(gomock.Any(), user).Return("", auth.ErrIncorrectData)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrIncorrectData.Error() + `"}`,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := authMocks.NewMockUsecase(ctrl)
			tt.mockBehavior(mockUsecase, tt.inputUser)

			handler := &Handler{
				usecase: mockUsecase,
				logger:  logger.SetupLogger(),
			}

			router := mux.NewRouter()
			router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.inputBody))

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestHandler_Register(t *testing.T) {
	type mockBehavior func(m *authMocks.MockUsecase, user *models.Users)

	testTable := []struct {
		name         string
		inputBody    string
		inputUser    *models.Users
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name: "successful registration moderator",
			inputBody: `{
				"email": "test@example.com",
				"password": "12345",
				"role": "moderator"
			}`,
			inputUser: &models.Users{
				Email:    "test@example.com",
				Password: "12345",
				Role:     "moderator",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {
				m.EXPECT().Register(gomock.Any(), user).Return("testToken", nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: `"testToken"`,
		},
		{
			name: "successful registration employee",
			inputBody: `{
				"email": "test@example.com",
				"password": "12345",
				"role": "employee"
			}`,
			inputUser: &models.Users{
				Email:    "test@example.com",
				Password: "12345",
				Role:     "employee",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {
				m.EXPECT().Register(gomock.Any(), user).Return("testToken", nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: `"testToken"`,
		},
		{
			name:         "empty body",
			inputBody:    `{}`,
			inputUser:    nil,
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrBadRequest.Error() + `"}`,
		},
		{
			name:         "invalid JSON",
			inputBody:    `{"email": 123}`,
			inputUser:    nil,
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "ошибка в чтении данных"}`,
		},
		{
			name: "missing required fields",
			inputBody: `{
				"email": "test@example.com"
			}`,
			inputUser:    nil,
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrBadRequest.Error() + `"}`,
		},
		{
			name: "user already exists",
			inputBody: `{
				"email": "exists@example.com",
				"password": "password123",
				"role": "employee"
			}`,
			inputUser: &models.Users{
				Email:    "exists@example.com",
				Password: "password123",
				Role:     "employee",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {
				m.EXPECT().Register(gomock.Any(), user).Return("", auth.ErrAlreadyExists)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrAlreadyExists.Error() + `"}`,
		},
		{
			name: "incorrect role",
			inputBody: `{
				"email": "test@example.com",
				"password": "password123",
				"role": "invalid_role"
			}`,
			inputUser: &models.Users{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "invalid_role",
			},
			mockBehavior: func(m *authMocks.MockUsecase, user *models.Users) {
				m.EXPECT().Register(gomock.Any(), user).Return("", auth.ErrIncorrectRole)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"msg": "` + auth.ErrIncorrectRole.Error() + `"}`,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := authMocks.NewMockUsecase(ctrl)
			tt.mockBehavior(mockUsecase, tt.inputUser)

			handler := &Handler{
				usecase: mockUsecase,
				logger:  logger.SetupLogger(),
			}

			router := mux.NewRouter()
			router.HandleFunc("/register", handler.Register).Methods(http.MethodPost)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.inputBody))

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
