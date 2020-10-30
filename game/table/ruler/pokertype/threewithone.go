package pokertype

import "ddz/game/poker"

// ThreeWithOne ThreeWithOne
type ThreeWithOne struct {
	threePokers []poker.IPoker
	onePokers   []poker.IPoker
}

// IsThreeWithOne IsThreeWithOne
func IsThreeWithOne(pokers []poker.IPoker) bool {
	if len(pokers) != 4 {
		return false
	}
	return (pokers[0].Value() != pokers[1].Value() &&
		pokers[1].Value() == pokers[2].Value() &&
		pokers[2].Value() == pokers[3].Value()) ||

		(pokers[4].Value() != pokers[3].Value() &&
			pokers[3].Value() == pokers[2].Value() &&
			pokers[2].Value() == pokers[1].Value())
}

// NewThreeWithOne NewThreeWithOne
func NewThreeWithOne(pokers []poker.IPoker) ThreeWithOne {
	threes := make([]poker.IPoker, 3)
	one := make([]poker.IPoker, 1)
	if pokers[0] == pokers[1] {
		one = pokers[0:0]
		threes = pokers[1:]
	} else {
		one = pokers[3:]
		threes = pokers[0:3]
	}
	return ThreeWithOne{
		threePokers: threes,
		onePokers:   one,
	}
}

// Type Type
func (p ThreeWithOne) Type() int {
	return TypeThreeWithOne
}

// Value Value
func (p ThreeWithOne) Value() int {
	return p.threePokers[0].Value() * 3
}

// Pokers Pokers
func (p ThreeWithOne) Pokers() []poker.IPoker {
	return append(p.threePokers, p.onePokers...)
}

// IsBoom IsBoom
func (p ThreeWithOne) IsBoom() bool {
	return false
}

// Length Length
func (p ThreeWithOne) Length() int {
	return 4
}
