package jwt

type Type int8

const (
	TypeManager Type = iota + 1
	TypeUser
	TypePreverify
)
