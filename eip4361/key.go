package eip4361

import "fmt"

type Key int

const (
	_ Key = iota
	Domain
	Address
	Statement
	URI
	Version
	ChainID
	Nonce
	IssuedAt
	ExpirationTime
	NotBefore
	RequestID
	Resources
)

//go:generate stringer -type=Key

func (key Key) err(err error) error {
	return &Err{
		err: err,
		Key: key,
	}
}

func (key Key) errorf(format string, a ...any) error {
	return key.err(fmt.Errorf(format, a...))
}

func (key Key) invalidFormat() error {
	return key.errorf("invalid format")
}
