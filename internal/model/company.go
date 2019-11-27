package model

type Company struct {
	ID       int
	Name     string
	Accounts []Account
	Version  int
	Audited
}
