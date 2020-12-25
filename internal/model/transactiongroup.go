package model

type TransactionGroup struct {
	ID          int
	Name        string
	Description string
	Version     int
	Audited
}
