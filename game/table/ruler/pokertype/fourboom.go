package pokertype

import "ddz/game/poker"

// FourBoom FourBoom
type FourBoom struct {
	pokers []poker.IPoker
}

// IsFourBoom IsFourBoom
func IsFourBoom(pokers []poker.IPoker) bool {
	if len(pokers) != 4 {
		return false
	}
	return pokers[0].Value() == pokers[1].Value() && pokers[1].Value() == pokers[2].Value() &&
		pokers[2].Value() == pokers[3].Value()
}

// NewFourBoom NewFourBoom
func NewFourBoom(pokers []poker.IPoker) FourBoom {
	return FourBoom{
		pokers: pokers,
	}
}

// Type Type
func (p FourBoom) Type() int {
	return TypeFourBoom
}

// Value Value
func (p FourBoom) Value() int {
	return p.pokers[0].Value()
}

// Pokers Pokers
func (p FourBoom) Pokers() []poker.IPoker {
	return p.pokers
}

// IsBoom IsBoom
func (p FourBoom) IsBoom() bool {
	return true
}

// Length Length
func (p FourBoom) Length() int {
	return len(p.pokers)
}
