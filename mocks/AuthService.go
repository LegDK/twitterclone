// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	twitterclone "twitterclone"

	mock "github.com/stretchr/testify/mock"
)

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, input
func (_m *AuthService) Login(ctx context.Context, input twitterclone.LoginInput) (twitterclone.AuthResponse, error) {
	ret := _m.Called(ctx, input)

	var r0 twitterclone.AuthResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, twitterclone.LoginInput) (twitterclone.AuthResponse, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, twitterclone.LoginInput) twitterclone.AuthResponse); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(twitterclone.AuthResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, twitterclone.LoginInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, input
func (_m *AuthService) Register(ctx context.Context, input twitterclone.RegisterInput) (twitterclone.AuthResponse, error) {
	ret := _m.Called(ctx, input)

	var r0 twitterclone.AuthResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, twitterclone.RegisterInput) (twitterclone.AuthResponse, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, twitterclone.RegisterInput) twitterclone.AuthResponse); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(twitterclone.AuthResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, twitterclone.RegisterInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthService interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthService creates a new instance of AuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthService(t mockConstructorTestingTNewAuthService) *AuthService {
	mock := &AuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
