package model

import (
	"database/sql/driver"
)

// YesNo maps a boolean value to a char column containing Y or N.
type YesNo struct {
	value byte
}

// Value returns the value to use for a query parameter.
func (yn *YesNo) Value() (driver.Value, error) {
	return yn.value, nil
}

// Scan populates a YesNo with the value from a query row.
func (yn *YesNo) Scan(src interface{}) error {
	yn.value = src.([]byte)[0]
	return nil
}

func (yn *YesNo) Get() bool {
	return yn.value == 'Y'
}

func (yn *YesNo) Set(value bool) {
	if value {
		yn.value = 'Y'
	} else {
		yn.value = 'N'
	}
}
