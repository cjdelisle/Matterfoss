// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/cjdelisle/matterfoss-server/v6/model"
	mock "github.com/stretchr/testify/mock"
)

// ProductNoticesStore is an autogenerated mock type for the ProductNoticesStore type
type ProductNoticesStore struct {
	mock.Mock
}

// Clear provides a mock function with given fields: notices
func (_m *ProductNoticesStore) Clear(notices []string) error {
	ret := _m.Called(notices)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(notices)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearOldNotices provides a mock function with given fields: currentNotices
func (_m *ProductNoticesStore) ClearOldNotices(currentNotices model.ProductNotices) error {
	ret := _m.Called(currentNotices)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.ProductNotices) error); ok {
		r0 = rf(currentNotices)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetViews provides a mock function with given fields: userID
func (_m *ProductNoticesStore) GetViews(userID string) ([]model.ProductNoticeViewState, error) {
	ret := _m.Called(userID)

	var r0 []model.ProductNoticeViewState
	if rf, ok := ret.Get(0).(func(string) []model.ProductNoticeViewState); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ProductNoticeViewState)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// View provides a mock function with given fields: userID, notices
func (_m *ProductNoticesStore) View(userID string, notices []string) error {
	ret := _m.Called(userID, notices)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(userID, notices)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
