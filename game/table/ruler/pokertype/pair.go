package pokertype

import "ddz/game/poker"

// Pair Pair
type Pair struct {
	pokers []poker.IPoker
}

// IsPair IsPair
func IsPair(pokers []poker.IPoker) bool {
	if len(pokers) == 2 && pokers[0].Value() == pokers[1].Value() {
		return true
	}
	return false
}

// NewPair NewPair
func NewPair(pokers []poker.IPoker) Pair {
	return Pair{
		pokers: pokers,
	}
}

// Type Type
func (p Pair) Type() int {
	return TypePair
}

// Value Value
func (p Pair) Value() int {
	return p.pokers[0].Value() * 2
}

// Pokers Pokers
func (p Pair) Pokers() []poker.IPoker {
	return p.pokers
}

// IsBoom IsBoom
func (p Pair) IsBoom() bool {
	return false
}

// Length Length
func (p Pair) Length() int {
	return len(p.pokers)
}
