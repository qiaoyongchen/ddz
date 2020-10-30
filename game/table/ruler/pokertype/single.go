package pokertype

import "ddz/game/poker"

// Single Single
type Single struct {
	pokers []poker.IPoker
}

// IsSingle IsSingle
func IsSingle(pokers []poker.IPoker) bool {
	if len(pokers) == 1 {
		return true
	}
	return false
}

// NewSingle NewSingle
func NewSingle(pokers []poker.IPoker) Single {
	return Single{
		pokers: pokers,
	}
}

// Type Type
func (p Single) Type() int {
	return TypeSingle
}

// Value Value
func (p Single) Value() int {
	return p.pokers[0].Value()
}

// Pokers Pokers
func (p Single) Pokers() []poker.IPoker {
	return p.pokers
}

// IsBoom IsBoom
func (p Single) IsBoom() bool {
	return false
}

// Length Length
func (p Single) Length() int {
	return len(p.pokers)
}
