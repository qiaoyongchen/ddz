package pokertype

import "ddz/game/poker"

// ThreeWithTwo ThreeWithTwo
type ThreeWithTwo struct {
	threePokers []poker.IPoker
	twoPokers   []poker.IPoker
}

// IsThreeWithTwo IsThreeWithTwo
func IsThreeWithTwo(pokers []poker.IPoker) bool {
	if len(pokers) != 5 {
		return false
	}
	return (pokers[0].Value() == pokers[1].Value() &&
		pokers[1].Value() != pokers[2].Value() &&
		pokers[2].Value() == pokers[3].Value() &&
		pokers[3].Value() == pokers[4].Value()) ||

		(pokers[0].Value() == pokers[1].Value() &&
			pokers[1].Value() == pokers[2].Value() &&
			pokers[2].Value() != pokers[3].Value() &&
			pokers[3].Value() == pokers[4].Value())
}

// NewThreeWithTwo NewThreeWithTwo
func NewThreeWithTwo(pokers []poker.IPoker) ThreeWithTwo {
	twos := make([]poker.IPoker, 2)
	threes := make([]poker.IPoker, 3)
	if pokers[1].Value() != pokers[2].Value() {
		twos = pokers[0:2]
		threes = pokers[2:]
	} else {
		twos = pokers[3:]
		threes = pokers[0:3]
	}
	return ThreeWithTwo{
		threePokers: threes,
		twoPokers:   twos,
	}
}

// Type Type
func (p ThreeWithTwo) Type() int {
	return TypeThreeWithTwo
}

// Value Value
func (p ThreeWithTwo) Value() int {
	return p.threePokers[0].Value() * 3
}

// Pokers Pokers
func (p ThreeWithTwo) Pokers() []poker.IPoker {
	return append(p.threePokers, p.twoPokers...)
}

// IsBoom IsBoom
func (p ThreeWithTwo) IsBoom() bool {
	return false
}

// Length Length
func (p ThreeWithTwo) Length() int {
	return 5
}
