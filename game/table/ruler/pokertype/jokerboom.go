package pokertype

import "ddz/game/poker"

// JokerBoom 王炸
type JokerBoom struct {
	pokers []poker.IPoker
}

// IsJokerBoom IsJokerBoom
func IsJokerBoom(pokers []poker.IPoker) bool {
	if len(pokers) != 2 {
		return false
	}
	return pokers[0].Value() == poker.VLittleJoker && pokers[1].Value() == poker.VBigJoker
}

// NewJokerBoom NewJokerBoom
func NewJokerBoom(pokers []poker.IPoker) JokerBoom {
	return JokerBoom{
		pokers: pokers,
	}
}

// Type Type
func (p JokerBoom) Type() int {
	return TypeJokerBoom
}

// Value Value
func (p JokerBoom) Value() int {
	pokersValue := 0
	for _, _p := range p.pokers {
		pokersValue += _p.Value()
	}
	return pokersValue * p.Type()
}

// Pokers Pokers
func (p JokerBoom) Pokers() []poker.IPoker {
	return p.pokers
}

// IsBoom IsBoom
func (p JokerBoom) IsBoom() bool {
	return true
}

// Length Length
func (p JokerBoom) Length() int {
	return len(p.pokers)
}
