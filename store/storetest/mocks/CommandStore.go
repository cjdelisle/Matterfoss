// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/cjdelisle/matterfoss-server/v6/model"
	mock "github.com/stretchr/testify/mock"
)

// CommandStore is an autogenerated mock type for the CommandStore type
type CommandStore struct {
	mock.Mock
}

// AnalyticsCommandCount provides a mock function with given fields: teamID
func (_m *CommandStore) AnalyticsCommandCount(teamID string) (int64, error) {
	ret := _m.Called(teamID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(teamID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(teamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: commandID, time
func (_m *CommandStore) Delete(commandID string, time int64) error {
	ret := _m.Called(commandID, time)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(commandID, time)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *CommandStore) Get(id string) (*model.Command, error) {
	ret := _m.Called(id)

	var r0 *model.Command
	if rf, ok := ret.Get(0).(func(string) *model.Command); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTeam provides a mock function with given fields: teamID
func (_m *CommandStore) GetByTeam(teamID string) ([]*model.Command, error) {
	ret := _m.Called(teamID)

	var r0 []*model.Command
	if rf, ok := ret.Get(0).(func(string) []*model.Command); ok {
		r0 = rf(teamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(teamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTrigger provides a mock function with given fields: teamID, trigger
func (_m *CommandStore) GetByTrigger(teamID string, trigger string) (*model.Command, error) {
	ret := _m.Called(teamID, trigger)

	var r0 *model.Command
	if rf, ok := ret.Get(0).(func(string, string) *model.Command); ok {
		r0 = rf(teamID, trigger)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(teamID, trigger)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDeleteByTeam provides a mock function with given fields: teamID
func (_m *CommandStore) PermanentDeleteByTeam(teamID string) error {
	ret := _m.Called(teamID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(teamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PermanentDeleteByUser provides a mock function with given fields: userID
func (_m *CommandStore) PermanentDeleteByUser(userID string) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: webhook
func (_m *CommandStore) Save(webhook *model.Command) (*model.Command, error) {
	ret := _m.Called(webhook)

	var r0 *model.Command
	if rf, ok := ret.Get(0).(func(*model.Command) *model.Command); ok {
		r0 = rf(webhook)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Command) error); ok {
		r1 = rf(webhook)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: hook
func (_m *CommandStore) Update(hook *model.Command) (*model.Command, error) {
	ret := _m.Called(hook)

	var r0 *model.Command
	if rf, ok := ret.Get(0).(func(*model.Command) *model.Command); ok {
		r0 = rf(hook)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Command) error); ok {
		r1 = rf(hook)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
