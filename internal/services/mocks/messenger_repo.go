// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	domain "messenger/internal/domain"

	mock "github.com/stretchr/testify/mock"

	models "messenger/internal/domain/models"

	uuid "github.com/google/uuid"
)

// MessengerRepo is an autogenerated mock type for the MessengerRepo type
type MessengerRepo struct {
	mock.Mock
}

// Add provides a mock function with given fields: message
func (_m *MessengerRepo) Add(message models.Message) (models.Message, error) {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 models.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(models.Message) (models.Message, error)); ok {
		return rf(message)
	}
	if rf, ok := ret.Get(0).(func(models.Message) models.Message); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Get(0).(models.Message)
	}

	if rf, ok := ret.Get(1).(func(models.Message) error); ok {
		r1 = rf(message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *MessengerRepo) Delete(id uuid.UUID) error {
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
func (_m *MessengerRepo) GetByChat(chatId uuid.UUID) ([]models.Message, error) {
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
func (_m *MessengerRepo) GetById(id uuid.UUID) (models.Message, error) {
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

// GetUserChats provides a mock function with given fields: userId
func (_m *MessengerRepo) GetUserChats(userId uuid.UUID) ([]uuid.UUID, error) {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserChats")
	}

	var r0 []uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) ([]uuid.UUID, error)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) []uuid.UUID); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: message
func (_m *MessengerRepo) Update(message domain.MessageUpdate) error {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.MessageUpdate) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMessengerRepo creates a new instance of MessengerRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMessengerRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MessengerRepo {
	mock := &MessengerRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
