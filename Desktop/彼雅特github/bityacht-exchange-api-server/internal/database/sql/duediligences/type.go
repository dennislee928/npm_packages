package duediligences

type Type int32

const (
	TypeCreateByIDV Type = iota + 1
	TypeManualSet
	TypeManualResend
)

type Bool int32

const (
	BoolUnknown Bool = iota
	BoolFalse
	BoolTrue
)

func (b *Bool) Set(val bool) {
	switch val {
	case false:
		*b = BoolFalse
	case true:
		*b = BoolTrue
	}
}
