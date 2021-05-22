package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_company_PtrTo(t *testing.T) {
	company := &Company{}
	tests := []struct {
		column string
		ptr    interface{}
	}{
		{column: "id", ptr: &company.ID},
		{column: "name", ptr: &company.Name},
		{column: "version", ptr: &company.Version},
		{column: "change_user", ptr: &company.ChangeUser},
		{column: "change_date", ptr: &company.ChangeDate},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s returns pointer to field", test.column), func(t *testing.T) {
			field := company.PtrTo(test.column)

			assert.Same(t, test.ptr, field)
		})
	}
}
