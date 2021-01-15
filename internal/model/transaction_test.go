package model

import (
	"reflect"
	"testing"
)

func Test_Transaction_GetRelatedDetailIDs(t *testing.T) {
	tests := []struct {
		description string
		relatedID   int64
		expectedIDs []int64
	}{
		{"does nothing if no related detail", -1, []int64{}},
		{"appends related detail ID", 1, []int64{1}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var relatedID *int64 = nil
			if test.relatedID > 0 {
				relatedID = &test.relatedID
			} else {
				relatedID = nil
			}
			details := append(make([]TransactionDetail, 0), TransactionDetail{RelatedDetailID: relatedID})
			tx := Transaction{Details: details}

			ids := tx.GetRelatedDetailIDs(make([]int64, 0))

			if !reflect.DeepEqual(ids, test.expectedIDs) {
				t.Errorf("Expected IDs: %v, got: %v", test.expectedIDs, ids)
			}
		})
	}
}
