package model

import "time"

type Audited struct {
	ChangeUser string
	ChangeDate *time.Time
}

type Versioned struct {
	Version int
}
