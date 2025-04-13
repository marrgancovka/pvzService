package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	"github.com/marrgancovka/pvzService/internal/services/pvz/mocks"
	"github.com/marrgancovka/pvzService/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestUsecase_CreatePvz(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	log := logger.SetupLogger()

	uc := Usecase{
		log:  log,
		repo: mockRepo,
	}

	validCity := models.CityMoscow
	invalidCity := models.City("Владивосток")
	testUUID := uuid.New()
	testTime := time.Now()

	tests := []struct {
		name        string
		input       *models.Pvz
		mockSetup   func()
		expected    *models.Pvz
		expectedErr error
	}{
		{
			name: "success with all fields",
			input: &models.Pvz{
				ID:               testUUID,
				RegistrationDate: testTime,
				City:             validCity,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreatePvz(gomock.Any(), gomock.Any()).
					Return(&models.Pvz{
						ID:               testUUID,
						City:             validCity,
						RegistrationDate: testTime,
					}, nil)
			},
			expected: &models.Pvz{
				ID:               testUUID,
				City:             validCity,
				RegistrationDate: testTime,
			},
			expectedErr: nil,
		},
		{
			name: "success with generated ID and date",
			input: &models.Pvz{
				City: validCity,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreatePvz(gomock.Any(), gomock.Any()).
					Return(&models.Pvz{
						ID:               testUUID,
						City:             validCity,
						RegistrationDate: testTime,
					}, nil)
			},
			expected: &models.Pvz{
				ID:               testUUID,
				RegistrationDate: testTime,
				City:             validCity,
			},
			expectedErr: nil,
		},
		{
			name: "invalid city",
			input: &models.Pvz{
				ID:               testUUID,
				RegistrationDate: testTime,
				City:             invalidCity,
			},
			mockSetup:   func() {},
			expected:    nil,
			expectedErr: pvz.ErrInaccessibleCity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.CreatePvz(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)

				if tt.input.ID == uuid.Nil {
					assert.NotEqual(t, uuid.Nil, result.ID)
				} else {
					assert.Equal(t, tt.input.ID, result.ID)
				}

				if tt.input.RegistrationDate.IsZero() {
					assert.False(t, result.RegistrationDate.IsZero())
				} else {
					assert.Equal(t, tt.input.RegistrationDate, result.RegistrationDate)
				}

				assert.Equal(t, tt.input.City, result.City)
			}
		})
	}
}

func TestUsecase_CreateReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	log := logger.SetupLogger()

	uc := Usecase{
		log:  log,
		repo: mockRepo,
	}

	testUUID := uuid.New()
	testPvzID := uuid.New()
	testTime := time.Now()

	tests := []struct {
		name         string
		input        *models.ReceptionRequest
		mockBehavior func()
		expected     *models.Reception
		expectedErr  error
	}{
		{
			name: "successful",
			input: &models.ReceptionRequest{
				PvzID: testPvzID,
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					CreateReception(gomock.Any(), gomock.Any()).
					Return(&models.Reception{
						ID:       testUUID,
						DateTime: testTime,
						PvzID:    testPvzID,
						Status:   models.StatusInProgress,
					}, nil)
			},
			expected: &models.Reception{
				ID:       testUUID,
				DateTime: testTime,
				PvzID:    testPvzID,
				Status:   models.StatusInProgress,
			},
			expectedErr: nil,
		},
		{
			name: "in progress reception already exists",
			input: &models.ReceptionRequest{
				PvzID: testPvzID,
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					CreateReception(gomock.Any(), gomock.Any()).
					Return(nil, pvz.ErrNoClosedReception)
			},
			expected:    nil,
			expectedErr: pvz.ErrNoClosedReception,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockBehavior != nil {
				tt.mockBehavior()
			}

			result, err := uc.CreateReception(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.False(t, result.DateTime.IsZero())
				assert.Equal(t, tt.input.PvzID, result.PvzID)
				assert.Equal(t, models.StatusInProgress, result.Status)
			}
		})
	}
}

func TestUsecase_CloseLastReceptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	log := logger.SetupLogger()

	uc := Usecase{
		log:  log,
		repo: mockRepo,
	}

	validPvzID := uuid.New()
	testReception := &models.Reception{
		ID:       uuid.New(),
		PvzID:    validPvzID,
		Status:   models.StatusClose,
		DateTime: time.Now(),
	}

	tests := []struct {
		name        string
		pvzId       uuid.UUID
		mockSetup   func()
		expected    *models.Reception
		expectedErr error
	}{
		{
			name:  "successful",
			pvzId: validPvzID,
			mockSetup: func() {
				mockRepo.EXPECT().
					CloseLastReceptions(gomock.Any(), validPvzID).
					Return(testReception, nil)
			},
			expected:    testReception,
			expectedErr: nil,
		},
		{
			name:  "no exists pvzID",
			pvzId: validPvzID,
			mockSetup: func() {
				mockRepo.EXPECT().
					CloseLastReceptions(gomock.Any(), validPvzID).
					Return(nil, pvz.ErrPvzNotExists)
			},
			expected:    nil,
			expectedErr: pvz.ErrPvzNotExists,
		},
		{
			name:  "no open reception",
			pvzId: validPvzID,
			mockSetup: func() {
				mockRepo.EXPECT().
					CloseLastReceptions(gomock.Any(), validPvzID).
					Return(nil, pvz.ErrNoOpenReception)
			},
			expected:    nil,
			expectedErr: pvz.ErrNoOpenReception,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.CloseLastReceptions(context.Background(), tt.pvzId)

			assert.Equal(t, tt.expected, result)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsecase_AddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	log := logger.SetupLogger()

	uc := Usecase{
		log:  log,
		repo: mockRepo,
	}

	validPvzID := uuid.New()
	validID := uuid.New()
	dateTime := time.Now()
	validProductType := models.TypeElectronics
	invalidProductType := models.ProductType("чай")

	tests := []struct {
		name        string
		input       *models.ProductRequest
		mockSetup   func()
		expected    *models.Product
		expectedErr error
	}{
		{
			name: "successful",
			input: &models.ProductRequest{
				Type:  validProductType,
				PvzID: validPvzID,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					AddProduct(gomock.Any(), gomock.Any(), validPvzID).
					Return(&models.Product{
						ID:          validID,
						DateTime:    dateTime,
						Type:        validProductType,
						ReceptionID: validID,
					}, nil)
			},
			expected: &models.Product{
				ID:          validID,
				DateTime:    dateTime,
				Type:        validProductType,
				ReceptionID: validID,
			},
			expectedErr: nil,
		},
		{
			name: "invalid product type",
			input: &models.ProductRequest{
				Type:  invalidProductType,
				PvzID: validPvzID,
			},
			mockSetup:   func() {},
			expected:    nil,
			expectedErr: pvz.ErrIncorrectProductType,
		},
		{
			name: "pvz is not exists",
			input: &models.ProductRequest{
				Type:  validProductType,
				PvzID: validPvzID,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					AddProduct(gomock.Any(), gomock.Any(), validPvzID).
					Return(nil, pvz.ErrPvzNotExists)
			},
			expected:    nil,
			expectedErr: pvz.ErrPvzNotExists,
		},
		{
			name: "no open reception",
			input: &models.ProductRequest{
				Type:  validProductType,
				PvzID: validPvzID,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					AddProduct(gomock.Any(), gomock.Any(), validPvzID).
					Return(nil, pvz.ErrNoOpenReception)
			},
			expected:    nil,
			expectedErr: pvz.ErrNoOpenReception,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.AddProduct(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, tt.input.Type, result.Type)
			}
		})
	}
}

func TestUsecase_DeleteLastProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	log := logger.SetupLogger()

	uc := Usecase{
		log:  log,
		repo: mockRepo,
	}

	validPvzID := uuid.New()

	tests := []struct {
		name        string
		pvzId       uuid.UUID
		mockSetup   func()
		expectedErr error
	}{
		{
			name:  "successful deletion",
			pvzId: validPvzID,
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteLastProduct(gomock.Any(), validPvzID).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "no products",
			pvzId: validPvzID,
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteLastProduct(gomock.Any(), validPvzID).
					Return(pvz.ErrNoProduct)
			},
			expectedErr: pvz.ErrNoProduct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.DeleteLastProduct(context.Background(), tt.pvzId)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
