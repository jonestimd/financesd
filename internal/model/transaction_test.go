package model

import (
	"reflect"
	"testing"
)

func Test_Transaction_GetRelatedDetailIDs(t *testing.T) {
	tests := []struct{
		description string
		relatedId 	int
		expectedIds []int
	} {
		{"does nothing if no related detail", -1, []int{}},
		{"appends related detail ID", 1, []int{1}},
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
			var relatedId *int = nil
			if test.relatedId > 0 {
				relatedId = &test.relatedId
			} else {
				relatedId = nil
			}
			details := append(make([]TransactionDetail, 0), TransactionDetail{RelatedDetailID: relatedId})
			tx := Transaction{Details: details}

			ids := tx.GetRelatedDetailIDs(make([]int, 0));

			if !reflect.DeepEqual(ids, test.expectedIds) {
				t.Errorf("Expected IDs: %v, got: %v", test.expectedIds, ids)
			}
		})
	}
}