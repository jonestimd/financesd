package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InputObject_GetInt(t *testing.T) {
	tests := []struct {
		name          string
		values        InputObject
		expectedValue interface{}
		expectedOk    bool
	}{
		{"returns nil, false for no value", map[string]interface{}{}, nil, false},
		{"returns nil, true for nil value", map[string]interface{}{"key": nil}, nil, true},
		{"returns value, true", map[string]interface{}{"key": 42}, int64(42), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := test.values.GetInt("key")

			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.expectedOk, ok)
		})
	}
}

func Test_InputObject_IntOrNull(t *testing.T) {
	tests := []struct {
		name          string
		values        InputObject
		expectedValue interface{}
	}{
		{"returns nil for no value", map[string]interface{}{}, nil},
		{"returns nil for nil value", map[string]interface{}{"key": nil}, nil},
		{"returns value", map[string]interface{}{"key": 42}, int64(42)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.values.IntOrNull("key")

			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func Test_InputObject_GetFloat(t *testing.T) {
	tests := []struct {
		name          string
		values        InputObject
		expectedValue interface{}
		expectedOk    bool
	}{
		{"returns nil, false for no value", map[string]interface{}{}, nil, false},
		{"returns nil, true for nil value", map[string]interface{}{"key": nil}, nil, true},
		{"returns value, true", map[string]interface{}{"key": 42.0}, 42.0, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := test.values.GetFloat("key")

			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.expectedOk, ok)
		})
	}
}

func Test_InputObject_FloatOrNull(t *testing.T) {
	tests := []struct {
		name          string
		values        InputObject
		expectedValue interface{}
	}{
		{"returns nil for no value", map[string]interface{}{}, nil},
		{"returns nil for nil value", map[string]interface{}{"key": nil}, nil},
		{"returns value", map[string]interface{}{"key": 42.0}, 42.0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.values.FloatOrNull("key")

			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func Test_InputObject_GetString(t *testing.T) {
	tests := []struct {
		name          string
		values        InputObject
		expectedValue interface{}
		expectedOk    bool
	}{
		{"returns nil, false for no value", map[string]interface{}{}, nil, false},
		{"returns nil, true for nil value", map[string]interface{}{"key": nil}, nil, true},
		{"returns value, true", map[string]interface{}{"key": "42"}, "42", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := test.values.GetString("key")

			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.expectedOk, ok)
		})
	}
}

func Test_InputObject_StringOrNull(t *testing.T) {
	tests := []struct {
		name          string
		values        InputObject
		expectedValue interface{}
	}{
		{"returns nil for no value", map[string]interface{}{}, nil},
		{"returns nil for nil value", map[string]interface{}{"key": nil}, nil},
		{"returns value", map[string]interface{}{"key": "42"}, "42"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.values.StringOrNull("key")

			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func Test_InputObject_GetVersionID(t *testing.T) {
	tests := []struct {
		name          string
		values        InputObject
		expectedValue *VersionID
		expectedErr   error
	}{
		{"returns nil for no values", map[string]interface{}{}, nil, nil},
		{"returns error for id without version", map[string]interface{}{"id": 42}, nil, errors.New("version is required for update/delete")},
		{"returns id and version", map[string]interface{}{"id": 42, "version": 96}, &VersionID{42, 96}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := test.values.GetVersionID()

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}
