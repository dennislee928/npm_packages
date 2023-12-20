package currencies

type Type int32

const (
	TypeFiat Type = iota + 1
	TypeCrypto
)
