package model

import "time"

type Transaction struct {
	ID              int
	Date            time.Time
	Memo            *string
	ReferenceNumber *string
	Cleared         *YesNo
	AccountId       int
	PayeeId         *int
	SecurityId      *int
	Version         int
	Audited
}
