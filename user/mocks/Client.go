// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"

	user "github.com/FirstVisit/go-truevault/user"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, username, password, attributes, groupIds, status, accessTokenNotValueAfter
func (_m *Client) Create(ctx context.Context, username string, password string, attributes string, groupIds []string, status user.Status, accessTokenNotValueAfter time.Time) (user.User, error) {
	ret := _m.Called(ctx, username, password, attributes, groupIds, status, accessTokenNotValueAfter)

	var r0 user.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, []string, user.Status, time.Time) user.User); ok {
		r0 = rf(ctx, username, password, attributes, groupIds, status, accessTokenNotValueAfter)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, []string, user.Status, time.Time) error); ok {
		r1 = rf(ctx, username, password, attributes, groupIds, status, accessTokenNotValueAfter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAPIKey provides a mock function with given fields: ctx, userID
func (_m *Client) CreateAPIKey(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAccessToken provides a mock function with given fields: ctx, userId, notValidAfter
func (_m *Client) CreateAccessToken(ctx context.Context, userId string, notValidAfter time.Time) error {
	ret := _m.Called(ctx, userId, notValidAfter)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) error); ok {
		r0 = rf(ctx, userId, notValidAfter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, userID
func (_m *Client) Delete(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, userId, full
func (_m *Client) Get(ctx context.Context, userId []string, full bool) ([]user.User, error) {
	ret := _m.Called(ctx, userId, full)

	var r0 []user.User
	if rf, ok := ret.Get(0).(func(context.Context, []string, bool) []user.User); ok {
		r0 = rf(ctx, userId, full)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, bool) error); ok {
		r1 = rf(ctx, userId, full)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, status, full
func (_m *Client) List(ctx context.Context, status user.Status, full bool) ([]user.User, error) {
	ret := _m.Called(ctx, status, full)

	var r0 []user.User
	if rf, ok := ret.Get(0).(func(context.Context, user.Status, bool) []user.User); ok {
		r0 = rf(ctx, status, full)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, user.Status, bool) error); ok {
		r1 = rf(ctx, status, full)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, userId, username, password, accessToken, accessTokenNotValueAfter, attributes, status
func (_m *Client) Update(ctx context.Context, userId string, username string, password string, accessToken string, accessTokenNotValueAfter time.Time, attributes string, status user.Status) (user.User, error) {
	ret := _m.Called(ctx, userId, username, password, accessToken, accessTokenNotValueAfter, attributes, status)

	var r0 user.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, time.Time, string, user.Status) user.User); ok {
		r0 = rf(ctx, userId, username, password, accessToken, accessTokenNotValueAfter, attributes, status)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, time.Time, string, user.Status) error); ok {
		r1 = rf(ctx, userId, username, password, accessToken, accessTokenNotValueAfter, attributes, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePassword provides a mock function with given fields: ctx, userId, password
func (_m *Client) UpdatePassword(ctx context.Context, userId string, password string) error {
	ret := _m.Called(ctx, userId, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, userId, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
