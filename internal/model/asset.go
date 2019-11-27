package model

type Asset struct {
	ID      int
	Name    string
	Type    string
	Scale   int
	Symbol  *string
	Version int
	Audited
}
