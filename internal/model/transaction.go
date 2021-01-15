package model

import "time"

type Transaction struct {
	ID              int64
	Date            time.Time
	Memo            *string
	ReferenceNumber *string
	Cleared         *YesNo
	AccountID       int64
	PayeeID         *int64
	SecurityID      *int64
	Details         []TransactionDetail
	Version         int
	Audited
}

func (tx *Transaction) GetRelatedDetailIDs(dest []int64) []int64 {
	for _, detail := range tx.Details {
		if detail.RelatedDetailID != nil {
			dest = append(dest, *detail.RelatedDetailID)
		}
	}
	return dest
}
