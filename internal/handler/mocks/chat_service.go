// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	domain "messenger/internal/domain"

	mock "github.com/stretchr/testify/mock"

	models "messenger/internal/domain/models"

	uuid "github.com/google/uuid"
)

// ChatService is an autogenerated mock type for the ChatService type
type ChatService struct {
	mock.Mock
}

// Add provides a mock function with given fields: chat
func (_m *ChatService) Add(chat domain.AddChat) (uuid.UUID, error) {
	ret := _m.Called(chat)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.AddChat) (uuid.UUID, error)); ok {
		return rf(chat)
	}
	if rf, ok := ret.Get(0).(func(domain.AddChat) uuid.UUID); ok {
		r0 = rf(chat)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.AddChat) error); ok {
		r1 = rf(chat)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddNewUser provides a mock function with given fields: chatId, userId
func (_m *ChatService) AddNewUser(chatId uuid.UUID, userId uuid.UUID) error {
	ret := _m.Called(chatId, userId)

	if len(ret) == 0 {
		panic("no return value specified for AddNewUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(chatId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: chatId
func (_m *ChatService) Delete(chatId uuid.UUID) error {
	ret := _m.Called(chatId)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(chatId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetInfoUserChats provides a mock function with given fields: userId, page, count
func (_m *ChatService) GetInfoUserChats(userId uuid.UUID, page uint, count uint) ([]domain.GetChat, error) {
	ret := _m.Called(userId, page, count)

	if len(ret) == 0 {
		panic("no return value specified for GetInfoUserChats")
	}

	var r0 []domain.GetChat
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uint, uint) ([]domain.GetChat, error)); ok {
		return rf(userId, page, count)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, uint, uint) []domain.GetChat); ok {
		r0 = rf(userId, page, count)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.GetChat)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, uint, uint) error); ok {
		r1 = rf(userId, page, count)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserChats provides a mock function with given fields: userId
func (_m *ChatService) GetUserChats(userId uuid.UUID) ([]uuid.UUID, error) {
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

// RemoveUser provides a mock function with given fields: chatId, userId
func (_m *ChatService) RemoveUser(chatId uuid.UUID, userId uuid.UUID) error {
	ret := _m.Called(chatId, userId)

	if len(ret) == 0 {
		panic("no return value specified for RemoveUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(chatId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: chat
func (_m *ChatService) Update(chat models.Chat) error {
	ret := _m.Called(chat)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Chat) error); ok {
		r0 = rf(chat)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewChatService creates a new instance of ChatService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChatService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChatService {
	mock := &ChatService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
