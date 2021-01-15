package model

type TransactionGroup struct {
	ID          int64
	Name        string
	Description string
	Version     int
	Audited
}
