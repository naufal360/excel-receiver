// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "excel-receiver/entity"

	mock "github.com/stretchr/testify/mock"
)

// TokenInterface is an autogenerated mock type for the TokenInterface type
type TokenInterface struct {
	mock.Mock
}

// GetTokenAuthentication provides a mock function with given fields: ctx, token
func (_m *TokenInterface) GetTokenAuthentication(ctx context.Context, token string) (*entity.Token, error) {
	ret := _m.Called(ctx, token)

	var r0 *entity.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.Token, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.Token); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTokenInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewTokenInterface creates a new instance of TokenInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenInterface(t mockConstructorTestingTNewTokenInterface) *TokenInterface {
	mock := &TokenInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
