package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_inputObject_getInt(t *testing.T) {
	tests := []struct {
		name          string
		values        inputObject
		expectedValue interface{}
		expectedOk    bool
	}{
		{"returns nil, false for no value", map[string]interface{}{}, nil, false},
		{"returns nil, true for nil value", map[string]interface{}{"key": nil}, nil, true},
		{"returns value, true", map[string]interface{}{"key": 42}, int64(42), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := test.values.getInt("key")

			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.expectedOk, ok)
		})
	}
}

func Test_inputObject_intOrNull(t *testing.T) {
	tests := []struct {
		name          string
		values        inputObject
		expectedValue interface{}
	}{
		{"returns nil for no value", map[string]interface{}{}, nil},
		{"returns nil for nil value", map[string]interface{}{"key": nil}, nil},
		{"returns value", map[string]interface{}{"key": 42}, int64(42)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.values.intOrNull("key")

			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func Test_inputObject_getFloat(t *testing.T) {
	tests := []struct {
		name          string
		values        inputObject
		expectedValue interface{}
		expectedOk    bool
	}{
		{"returns nil, false for no value", map[string]interface{}{}, nil, false},
		{"returns nil, true for nil value", map[string]interface{}{"key": nil}, nil, true},
		{"returns value, true", map[string]interface{}{"key": 42.0}, 42.0, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := test.values.getFloat("key")

			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.expectedOk, ok)
		})
	}
}

func Test_inputObject_floatOrNull(t *testing.T) {
	tests := []struct {
		name          string
		values        inputObject
		expectedValue interface{}
	}{
		{"returns nil for no value", map[string]interface{}{}, nil},
		{"returns nil for nil value", map[string]interface{}{"key": nil}, nil},
		{"returns value", map[string]interface{}{"key": 42.0}, 42.0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.values.floatOrNull("key")

			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func Test_inputObject_getString(t *testing.T) {
	tests := []struct {
		name          string
		values        inputObject
		expectedValue interface{}
		expectedOk    bool
	}{
		{"returns nil, false for no value", map[string]interface{}{}, nil, false},
		{"returns nil, true for nil value", map[string]interface{}{"key": nil}, nil, true},
		{"returns value, true", map[string]interface{}{"key": "42"}, "42", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := test.values.getString("key")

			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.expectedOk, ok)
		})
	}
}

func Test_inputObject_stringOrNull(t *testing.T) {
	tests := []struct {
		name          string
		values        inputObject
		expectedValue interface{}
	}{
		{"returns nil for no value", map[string]interface{}{}, nil},
		{"returns nil for nil value", map[string]interface{}{"key": nil}, nil},
		{"returns value", map[string]interface{}{"key": "42"}, "42"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.values.stringOrNull("key")

			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func Test_inputObject_getVersionID(t *testing.T) {
	tests := []struct {
		name          string
		values        inputObject
		expectedValue *versionID
		expectedErr   error
	}{
		{"returns nil for no values", map[string]interface{}{}, nil, nil},
		{"returns error for id without version", map[string]interface{}{"id": 42}, nil, errors.New("version is required for update/delete")},
		{"returns id and version", map[string]interface{}{"id": 42, "version": 96}, &versionID{42, 96}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := test.values.getVersionID()

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}