package model

type Payee struct {
	ID      int64
	Name    string
	Version int
	Audited
}
