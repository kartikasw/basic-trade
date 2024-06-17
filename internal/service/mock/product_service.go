// Code generated by MockGen. DO NOT EDIT.
// Source: basic-trade/internal/service (interfaces: ProductService)
//
// Generated by this command:
//
//	mockgen -package mockService -destination internal/service/mock/product_service.go basic-trade/internal/service ProductService
//

// Package mockService is a generated GoMock package.
package mockService

import (
	entity "basic-trade/internal/entity"
	context "context"
	multipart "mime/multipart"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockProductService is a mock of ProductService interface.
type MockProductService struct {
	ctrl     *gomock.Controller
	recorder *MockProductServiceMockRecorder
}

// MockProductServiceMockRecorder is the mock recorder for MockProductService.
type MockProductServiceMockRecorder struct {
	mock *MockProductService
}

// NewMockProductService creates a new mock instance.
func NewMockProductService(ctrl *gomock.Controller) *MockProductService {
	mock := &MockProductService{ctrl: ctrl}
	mock.recorder = &MockProductServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductService) EXPECT() *MockProductServiceMockRecorder {
	return m.recorder
}

// CreateProduct mocks base method.
func (m *MockProductService) CreateProduct(arg0 context.Context, arg1 entity.Product, arg2 uuid.UUID, arg3 *multipart.FileHeader) (entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockProductServiceMockRecorder) CreateProduct(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockProductService)(nil).CreateProduct), arg0, arg1, arg2, arg3)
}

// DeleteProduct mocks base method.
func (m *MockProductService) DeleteProduct(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockProductServiceMockRecorder) DeleteProduct(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockProductService)(nil).DeleteProduct), arg0, arg1, arg2)
}

// GetAllProducts mocks base method.
func (m *MockProductService) GetAllProducts(arg0, arg1 int32) ([]entity.ProductView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProducts", arg0, arg1)
	ret0, _ := ret[0].([]entity.ProductView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProducts indicates an expected call of GetAllProducts.
func (mr *MockProductServiceMockRecorder) GetAllProducts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProducts", reflect.TypeOf((*MockProductService)(nil).GetAllProducts), arg0, arg1)
}

// GetProduct mocks base method.
func (m *MockProductService) GetProduct(arg0 uuid.UUID) (entity.ProductView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", arg0)
	ret0, _ := ret[0].(entity.ProductView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockProductServiceMockRecorder) GetProduct(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockProductService)(nil).GetProduct), arg0)
}

// SearchProducts mocks base method.
func (m *MockProductService) SearchProducts(arg0 string, arg1, arg2 int32) ([]entity.ProductView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchProducts", arg0, arg1, arg2)
	ret0, _ := ret[0].([]entity.ProductView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchProducts indicates an expected call of SearchProducts.
func (mr *MockProductServiceMockRecorder) SearchProducts(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchProducts", reflect.TypeOf((*MockProductService)(nil).SearchProducts), arg0, arg1, arg2)
}

// UpdateProduct mocks base method.
func (m *MockProductService) UpdateProduct(arg0 context.Context, arg1 entity.Product, arg2 uuid.UUID, arg3 *multipart.FileHeader) (entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockProductServiceMockRecorder) UpdateProduct(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockProductService)(nil).UpdateProduct), arg0, arg1, arg2, arg3)
}