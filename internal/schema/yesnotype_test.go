package schema

import (
	"testing"

	"github.com/jonestimd/financesd/internal/database/table"
)

func Test_yesnotype_Serialize(t *testing.T) {
	y := table.YesNo('Y')
	n := table.YesNo('N')
	tests := []struct {
		name   string
		value  interface{}
		result interface{}
	}{
		{"returns true for Y", y, true},
		{"returns true for pointer to Y", &y, true},
		{"returns false for N", n, false},
		{"returns true for pointer to N", &n, false},
		{"returns nil for non-string", 123, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := yesNoType.Serialize(test.value)

			if result != test.result {
				t.Errorf("Expected: %v, got %v", test.result, result)
			}
		})
	}
}
