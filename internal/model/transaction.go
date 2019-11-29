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
	Details         []TransactionDetail
	Version         int
	Audited
}

func (tx *Transaction) GetRelatedDetailIDs(dest []int) []int {
	for _, detail := range tx.Details {
		if detail.RelatedDetailID != nil {
			dest = append(dest, *detail.RelatedDetailID)
		}
	}
	return dest
}
