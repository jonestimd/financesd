package model

type Asset struct {
	ID      int64
	Name    string
	Type    string
	Scale   int
	Symbol  *string
	Version int
	Audited
}
