package pokertype

import "ddz/game/poker"

// JokerBoom 王炸
type JokerBoom struct {
	pokers []poker.IPoker
}

// NewJokerBoom NewJokerBoom
func NewJokerBoom(pokers []poker.IPoker) JokerBoom {
	if pokers[0].Value() > pokers[1].Value() {
		pokers[0], pokers[1] = pokers[1], pokers[0]
	}
	return JokerBoom{
		pokers: pokers,
	}
}

// Type Type
func (p JokerBoom) Type() int {
	return 1000
}

// Value Value
func (p JokerBoom) Value() int {
	pokersValue := 0
	for _, p := range p.pokers {
		pokersValue += p.Value()
	}
	return pokersValue * p.Type()
}

// Pokers Pokers
func (p JokerBoom) Pokers() []poker.IPoker {
	return p.pokers
}
