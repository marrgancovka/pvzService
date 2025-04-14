// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/pvz/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=internal/services/pvz/interfaces.go -destination=internal/services/pvz/mocks/mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	models "github.com/marrgancovka/pvzService/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
	isgomock struct{}
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockUsecase) AddProduct(ctx context.Context, product *models.ProductRequest) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, product)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockUsecaseMockRecorder) AddProduct(ctx, product any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockUsecase)(nil).AddProduct), ctx, product)
}

// CloseLastReceptions mocks base method.
func (m *MockUsecase) CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseLastReceptions", ctx, pvzId)
	ret0, _ := ret[0].(*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseLastReceptions indicates an expected call of CloseLastReceptions.
func (mr *MockUsecaseMockRecorder) CloseLastReceptions(ctx, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseLastReceptions", reflect.TypeOf((*MockUsecase)(nil).CloseLastReceptions), ctx, pvzId)
}

// CreatePvz mocks base method.
func (m *MockUsecase) CreatePvz(ctx context.Context, pvzData *models.Pvz) (*models.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePvz", ctx, pvzData)
	ret0, _ := ret[0].(*models.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePvz indicates an expected call of CreatePvz.
func (mr *MockUsecaseMockRecorder) CreatePvz(ctx, pvzData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePvz", reflect.TypeOf((*MockUsecase)(nil).CreatePvz), ctx, pvzData)
}

// CreateReception mocks base method.
func (m *MockUsecase) CreateReception(ctx context.Context, receptionData *models.ReceptionRequest) (*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReception", ctx, receptionData)
	ret0, _ := ret[0].(*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReception indicates an expected call of CreateReception.
func (mr *MockUsecaseMockRecorder) CreateReception(ctx, receptionData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReception", reflect.TypeOf((*MockUsecase)(nil).CreateReception), ctx, receptionData)
}

// DeleteLastProduct mocks base method.
func (m *MockUsecase) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", ctx, pvzId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockUsecaseMockRecorder) DeleteLastProduct(ctx, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockUsecase)(nil).DeleteLastProduct), ctx, pvzId)
}

// GetPvz mocks base method.
func (m *MockUsecase) GetPvz(ctx context.Context, startDate, endDate time.Time, limit, page uint64) ([]*models.PvzWithReceptions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvz", ctx, startDate, endDate, limit, page)
	ret0, _ := ret[0].([]*models.PvzWithReceptions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvz indicates an expected call of GetPvz.
func (mr *MockUsecaseMockRecorder) GetPvz(ctx, startDate, endDate, limit, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvz", reflect.TypeOf((*MockUsecase)(nil).GetPvz), ctx, startDate, endDate, limit, page)
}

// GetPvzList mocks base method.
func (m *MockUsecase) GetPvzList(ctx context.Context) ([]*models.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvzList", ctx)
	ret0, _ := ret[0].([]*models.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvzList indicates an expected call of GetPvzList.
func (mr *MockUsecaseMockRecorder) GetPvzList(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvzList", reflect.TypeOf((*MockUsecase)(nil).GetPvzList), ctx)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockRepository) AddProduct(ctx context.Context, product *models.Product, pvzID uuid.UUID) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, product, pvzID)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockRepositoryMockRecorder) AddProduct(ctx, product, pvzID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockRepository)(nil).AddProduct), ctx, product, pvzID)
}

// CloseLastReceptions mocks base method.
func (m *MockRepository) CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseLastReceptions", ctx, pvzId)
	ret0, _ := ret[0].(*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseLastReceptions indicates an expected call of CloseLastReceptions.
func (mr *MockRepositoryMockRecorder) CloseLastReceptions(ctx, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseLastReceptions", reflect.TypeOf((*MockRepository)(nil).CloseLastReceptions), ctx, pvzId)
}

// CreatePvz mocks base method.
func (m *MockRepository) CreatePvz(ctx context.Context, pvzData *models.Pvz) (*models.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePvz", ctx, pvzData)
	ret0, _ := ret[0].(*models.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePvz indicates an expected call of CreatePvz.
func (mr *MockRepositoryMockRecorder) CreatePvz(ctx, pvzData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePvz", reflect.TypeOf((*MockRepository)(nil).CreatePvz), ctx, pvzData)
}

// CreateReception mocks base method.
func (m *MockRepository) CreateReception(ctx context.Context, receptionData *models.Reception) (*models.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReception", ctx, receptionData)
	ret0, _ := ret[0].(*models.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReception indicates an expected call of CreateReception.
func (mr *MockRepositoryMockRecorder) CreateReception(ctx, receptionData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReception", reflect.TypeOf((*MockRepository)(nil).CreateReception), ctx, receptionData)
}

// DeleteLastProduct mocks base method.
func (m *MockRepository) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", ctx, pvzId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockRepositoryMockRecorder) DeleteLastProduct(ctx, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockRepository)(nil).DeleteLastProduct), ctx, pvzId)
}

// GetPvz mocks base method.
func (m *MockRepository) GetPvz(ctx context.Context, startDate, endDate time.Time, limit, page uint64) ([]*models.PvzWithReceptions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvz", ctx, startDate, endDate, limit, page)
	ret0, _ := ret[0].([]*models.PvzWithReceptions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvz indicates an expected call of GetPvz.
func (mr *MockRepositoryMockRecorder) GetPvz(ctx, startDate, endDate, limit, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvz", reflect.TypeOf((*MockRepository)(nil).GetPvz), ctx, startDate, endDate, limit, page)
}

// GetPvzList mocks base method.
func (m *MockRepository) GetPvzList(ctx context.Context) ([]*models.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvzList", ctx)
	ret0, _ := ret[0].([]*models.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvzList indicates an expected call of GetPvzList.
func (mr *MockRepositoryMockRecorder) GetPvzList(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvzList", reflect.TypeOf((*MockRepository)(nil).GetPvzList), ctx)
}
