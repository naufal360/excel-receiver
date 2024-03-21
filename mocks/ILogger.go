// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	logrus "github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"

	provider "excel-receiver/provider"
)

// ILogger is an autogenerated mock type for the ILogger type
type ILogger struct {
	mock.Mock
}

// Debugf provides a mock function with given fields: logType, format, args
func (_m *ILogger) Debugf(logType provider.LogType, format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, logType, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// ErrorWithFields provides a mock function with given fields: logType, fields, format, args
func (_m *ILogger) ErrorWithFields(logType provider.LogType, fields map[string]interface{}, format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, logType, fields, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Errorf provides a mock function with given fields: logType, format, args
func (_m *ILogger) Errorf(logType provider.LogType, format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, logType, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// InfoWithFields provides a mock function with given fields: logType, fields, format, args
func (_m *ILogger) InfoWithFields(logType provider.LogType, fields map[string]interface{}, format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, logType, fields, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Infof provides a mock function with given fields: logType, format, args
func (_m *ILogger) Infof(logType provider.LogType, format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, logType, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// WithFields provides a mock function with given fields: logType, fields
func (_m *ILogger) WithFields(logType provider.LogType, fields logrus.Fields) *logrus.Entry {
	ret := _m.Called(logType, fields)

	var r0 *logrus.Entry
	if rf, ok := ret.Get(0).(func(provider.LogType, logrus.Fields) *logrus.Entry); ok {
		r0 = rf(logType, fields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*logrus.Entry)
		}
	}

	return r0
}

type mockConstructorTestingTNewILogger interface {
	mock.TestingT
	Cleanup(func())
}

// NewILogger creates a new instance of ILogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewILogger(t mockConstructorTestingTNewILogger) *ILogger {
	mock := &ILogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
