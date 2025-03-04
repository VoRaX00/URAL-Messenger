// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	models "messenger/internal/domain/models"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0
func (_m *Repository) Add(_a0 models.Message) (models.Message, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 models.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(models.Message) (models.Message, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(models.Message) models.Message); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(models.Message)
	}

	if rf, ok := ret.Get(1).(func(models.Message) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *Repository) Delete(id uuid.UUID) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByChat provides a mock function with given fields: chatId
func (_m *Repository) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
	ret := _m.Called(chatId)

	if len(ret) == 0 {
		panic("no return value specified for GetByChat")
	}

	var r0 []models.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) ([]models.Message, error)); ok {
		return rf(chatId)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) []models.Message); ok {
		r0 = rf(chatId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(chatId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *Repository) GetById(id uuid.UUID) (models.Message, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 models.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (models.Message, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) models.Message); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(models.Message)
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *Repository) Update(_a0 models.Message) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Message) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
