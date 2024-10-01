// Code generated by mockery v2.46.1. DO NOT EDIT.

package redis

import mock "github.com/stretchr/testify/mock"

// MockPoolOption is an autogenerated mock type for the PoolOption type
type MockPoolOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: p
func (_m *MockPoolOption) Execute(p Pool) {
	_m.Called(p)
}

// NewMockPoolOption creates a new instance of MockPoolOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPoolOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPoolOption {
	mock := &MockPoolOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
