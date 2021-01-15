package graphql

import (
	"testing"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func Test_datetype_Serialize(t *testing.T) {
	now := time.Now()
	nowString := now.Format(dateFormat)
	tests := []struct {
		name   string
		value  interface{}
		result *string
	}{
		{"formats Time", now, &nowString},
		{"formats Time pointer", &now, &nowString},
		{"returns nil for non-Time", nil, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := dateType.Serialize(test.value)

			if test.result != nil {
				if result != *test.result {
					t.Errorf("Expected %v, got %v", test.result, result)
				}
			} else if result != nil {
				t.Error("Expected nil")
			}
		})
	}
}

func Test_datetype_ParseValue(t *testing.T) {
	todayString := time.Now().Format(dateFormat)
	today, _ := time.Parse(dateFormat, todayString)
	tests := []struct {
		name   string
		value  interface{}
		result *time.Time
	}{
		{"parses string", todayString, &today},
		{"parses string pointer", &todayString, &today},
		{"returns nil for non-string", nil, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := dateType.ParseValue(test.value)

			if test.result != nil {
				if result != *test.result {
					t.Errorf("Expected %v, got %v", test.result, result)
				}
			} else if result != nil {
				t.Error("Expected nil")
			}
		})
	}
}

func Test_datetype_ParseLiteral(t *testing.T) {
	todayString := time.Now().Format(dateFormat)
	today, _ := time.Parse(dateFormat, todayString)
	tests := []struct {
		name   string
		value  ast.Value
		result *time.Time
	}{
		{"parses string", &ast.StringValue{Kind: kinds.StringValue, Value: todayString}, &today},
		{"returns nil for invalid date", &ast.StringValue{Kind: kinds.StringValue, Value: "not a date"}, nil},
		{"returns nil for non-string", nil, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := dateType.ParseLiteral(test.value)

			if test.result != nil {
				if result != *test.result {
					t.Errorf("Expected %v, got %v", test.result, result)
				}
			} else if result != nil {
				t.Error("Expected nil")
			}
		})
	}
}
