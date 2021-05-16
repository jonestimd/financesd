package database

import (
	"database/sql/driver"
	"fmt"
)

// YesNo maps a boolean value to a char column containing Y or N.
type YesNo byte

// Value implements the sql.Valuer interface.
func (yn *YesNo) Value() (driver.Value, error) {
	return []byte{byte(*yn)}, nil
}

// Scan implements the sql.Scanner interface.
func (yn *YesNo) Scan(value interface{}) error {
	if b, ok := value.([]byte); ok && len(b) > 0 {
		*yn = YesNo(b[0])
		return nil
	}
	return fmt.Errorf("invalid value for YesNo: %v", value)
}

// Get returns the boolean value
func (yn *YesNo) Get() bool {
	return *yn == 'Y'
}

// Set sets the boolean value
func (yn *YesNo) Set(value bool) {
	if value {
		*yn = 'Y'
	} else {
		*yn = 'N'
	}
}
