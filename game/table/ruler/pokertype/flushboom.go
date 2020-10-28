package pokertype

import "ddz/game/poker"

// FlushBoom FlushBoom
type FlushBoom struct {
	pokers []poker.IPoker
}

// NewFlushBoom NewFlushBoom
func NewFlushBoom(pokers []poker.IPoker) FlushBoom {
	return FlushBoom{
		pokers: pokers,
	}
}

// Type Type
func (p FlushBoom) Type() int {
	return TypeFlushBoom
}

// Value Value
func (p FlushBoom) Value() int {
	return (p.pokers[0].Value() + 10) * p.pokers[0].Type()
}

// Pokers Pokers
func (p FlushBoom) Pokers() []poker.IPoker {
	return p.pokers
}

// IsBoom IsBoom
func (p FlushBoom) IsBoom() bool {
	return true
}
