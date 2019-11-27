package model

type Company struct {
	ID       int
	Name     string
	Accounts []Account
	Versioned
	Audited
}
