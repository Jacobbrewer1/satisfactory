// Code generated by mockery v2.46.0. DO NOT EDIT.

package vault

import (
	context "context"

	api "github.com/hashicorp/vault/api"

	mock "github.com/stretchr/testify/mock"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// Client provides a mock function with given fields:
func (_m *MockClient) Client() *api.Client {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Client")
	}

	var r0 *api.Client
	if rf, ok := ret.Get(0).(func() *api.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Client)
		}
	}

	return r0
}

// GetKvSecretV2 provides a mock function with given fields: ctx, name
func (_m *MockClient) GetKvSecretV2(ctx context.Context, name string) (*api.KVSecret, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetKvSecretV2")
	}

	var r0 *api.KVSecret
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*api.KVSecret, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *api.KVSecret); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.KVSecret)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecret provides a mock function with given fields: ctx, path
func (_m *MockClient) GetSecret(ctx context.Context, path string) (*api.Secret, error) {
	ret := _m.Called(ctx, path)

	if len(ret) == 0 {
		panic("no return value specified for GetSecret")
	}

	var r0 *api.Secret
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*api.Secret, error)); ok {
		return rf(ctx, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *api.Secret); ok {
		r0 = rf(ctx, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Secret)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransitDecrypt provides a mock function with given fields: ctx, data
func (_m *MockClient) TransitDecrypt(ctx context.Context, data string) (string, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for TransitDecrypt")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransitEncrypt provides a mock function with given fields: ctx, data
func (_m *MockClient) TransitEncrypt(ctx context.Context, data string) (*api.Secret, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for TransitEncrypt")
	}

	var r0 *api.Secret
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*api.Secret, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *api.Secret); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Secret)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockClient creates a new instance of MockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClient {
	mock := &MockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
