// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/nicewook/gocore/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// ProductRepository is an autogenerated mock type for the ProductRepository type
type ProductRepository struct {
	mock.Mock
}

type ProductRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ProductRepository) EXPECT() *ProductRepository_Expecter {
	return &ProductRepository_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: ctx
func (_m *ProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.Product, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Product); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type ProductRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ProductRepository_Expecter) GetAll(ctx interface{}) *ProductRepository_GetAll_Call {
	return &ProductRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *ProductRepository_GetAll_Call) Run(run func(ctx context.Context)) *ProductRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ProductRepository_GetAll_Call) Return(_a0 []domain.Product, _a1 error) *ProductRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_GetAll_Call) RunAndReturn(run func(context.Context) ([]domain.Product, error)) *ProductRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*domain.Product, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.Product); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type ProductRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *ProductRepository_Expecter) GetByID(ctx interface{}, id interface{}) *ProductRepository_GetByID_Call {
	return &ProductRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *ProductRepository_GetByID_Call) Run(run func(ctx context.Context, id int64)) *ProductRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *ProductRepository_GetByID_Call) Return(_a0 *domain.Product, _a1 error) *ProductRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_GetByID_Call) RunAndReturn(run func(context.Context, int64) (*domain.Product, error)) *ProductRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: ctx, product
func (_m *ProductRepository) Save(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	ret := _m.Called(ctx, product)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Product) (*domain.Product, error)); ok {
		return rf(ctx, product)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Product) *domain.Product); ok {
		r0 = rf(ctx, product)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Product) error); ok {
		r1 = rf(ctx, product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type ProductRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - product *domain.Product
func (_e *ProductRepository_Expecter) Save(ctx interface{}, product interface{}) *ProductRepository_Save_Call {
	return &ProductRepository_Save_Call{Call: _e.mock.On("Save", ctx, product)}
}

func (_c *ProductRepository_Save_Call) Run(run func(ctx context.Context, product *domain.Product)) *ProductRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Product))
	})
	return _c
}

func (_c *ProductRepository_Save_Call) Return(_a0 *domain.Product, _a1 error) *ProductRepository_Save_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductRepository_Save_Call) RunAndReturn(run func(context.Context, *domain.Product) (*domain.Product, error)) *ProductRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewProductRepository creates a new instance of ProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductRepository {
	mock := &ProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
