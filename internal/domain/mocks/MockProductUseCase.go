// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	domain "github.com/nicewook/gocore/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// ProductUseCase is an autogenerated mock type for the ProductUseCase type
type ProductUseCase struct {
	mock.Mock
}

type ProductUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *ProductUseCase) EXPECT() *ProductUseCase_Expecter {
	return &ProductUseCase_Expecter{mock: &_m.Mock}
}

// CreateProduct provides a mock function with given fields: product
func (_m *ProductUseCase) CreateProduct(product *domain.Product) (*domain.Product, error) {
	ret := _m.Called(product)

	if len(ret) == 0 {
		panic("no return value specified for CreateProduct")
	}

	var r0 *domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(*domain.Product) (*domain.Product, error)); ok {
		return rf(product)
	}
	if rf, ok := ret.Get(0).(func(*domain.Product) *domain.Product); ok {
		r0 = rf(product)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(*domain.Product) error); ok {
		r1 = rf(product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductUseCase_CreateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateProduct'
type ProductUseCase_CreateProduct_Call struct {
	*mock.Call
}

// CreateProduct is a helper method to define mock.On call
//   - product *domain.Product
func (_e *ProductUseCase_Expecter) CreateProduct(product interface{}) *ProductUseCase_CreateProduct_Call {
	return &ProductUseCase_CreateProduct_Call{Call: _e.mock.On("CreateProduct", product)}
}

func (_c *ProductUseCase_CreateProduct_Call) Run(run func(product *domain.Product)) *ProductUseCase_CreateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.Product))
	})
	return _c
}

func (_c *ProductUseCase_CreateProduct_Call) Return(_a0 *domain.Product, _a1 error) *ProductUseCase_CreateProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductUseCase_CreateProduct_Call) RunAndReturn(run func(*domain.Product) (*domain.Product, error)) *ProductUseCase_CreateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with no fields
func (_m *ProductUseCase) GetAll() ([]domain.Product, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.Product, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.Product); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductUseCase_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type ProductUseCase_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *ProductUseCase_Expecter) GetAll() *ProductUseCase_GetAll_Call {
	return &ProductUseCase_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *ProductUseCase_GetAll_Call) Run(run func()) *ProductUseCase_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ProductUseCase_GetAll_Call) Return(_a0 []domain.Product, _a1 error) *ProductUseCase_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductUseCase_GetAll_Call) RunAndReturn(run func() ([]domain.Product, error)) *ProductUseCase_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: id
func (_m *ProductUseCase) GetByID(id int64) (*domain.Product, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*domain.Product, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *domain.Product); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductUseCase_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type ProductUseCase_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - id int64
func (_e *ProductUseCase_Expecter) GetByID(id interface{}) *ProductUseCase_GetByID_Call {
	return &ProductUseCase_GetByID_Call{Call: _e.mock.On("GetByID", id)}
}

func (_c *ProductUseCase_GetByID_Call) Run(run func(id int64)) *ProductUseCase_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *ProductUseCase_GetByID_Call) Return(_a0 *domain.Product, _a1 error) *ProductUseCase_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductUseCase_GetByID_Call) RunAndReturn(run func(int64) (*domain.Product, error)) *ProductUseCase_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewProductUseCase creates a new instance of ProductUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductUseCase {
	mock := &ProductUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
