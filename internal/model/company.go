package model

type Company struct {
	ID       int64
	Name     string
	Accounts []Account
	Version  int
	Audited
}
