// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/nicewook/gocore/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserUseCase is an autogenerated mock type for the UserUseCase type
type UserUseCase struct {
	mock.Mock
}

type UserUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *UserUseCase) EXPECT() *UserUseCase_Expecter {
	return &UserUseCase_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: ctx
func (_m *UserUseCase) GetAll(ctx context.Context) ([]domain.User, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserUseCase_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type UserUseCase_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserUseCase_Expecter) GetAll(ctx interface{}) *UserUseCase_GetAll_Call {
	return &UserUseCase_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *UserUseCase_GetAll_Call) Run(run func(ctx context.Context)) *UserUseCase_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserUseCase_GetAll_Call) Return(_a0 []domain.User, _a1 error) *UserUseCase_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserUseCase_GetAll_Call) RunAndReturn(run func(context.Context) ([]domain.User, error)) *UserUseCase_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *UserUseCase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*domain.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserUseCase_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type UserUseCase_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *UserUseCase_Expecter) GetByID(ctx interface{}, id interface{}) *UserUseCase_GetByID_Call {
	return &UserUseCase_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *UserUseCase_GetByID_Call) Run(run func(ctx context.Context, id int64)) *UserUseCase_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *UserUseCase_GetByID_Call) Return(_a0 *domain.User, _a1 error) *UserUseCase_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserUseCase_GetByID_Call) RunAndReturn(run func(context.Context, int64) (*domain.User, error)) *UserUseCase_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserUseCase creates a new instance of UserUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserUseCase {
	mock := &UserUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
