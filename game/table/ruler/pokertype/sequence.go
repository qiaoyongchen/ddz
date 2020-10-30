package pokertype

import "ddz/game/poker"

// Sequence Sequence
type Sequence struct {
	pokers []poker.IPoker
}

// IsSequence IsSequence
func IsSequence(pokers []poker.IPoker) bool {
	if len(pokers) < 5 {
		return false
	}
	types := map[int]int{
		poker.TypeDiamond:     0,
		poker.TypeClub:        0,
		poker.TypeHeart:       0,
		poker.TypeLittleJoker: 0,
		poker.TypeBigJoker:    0,
	}
	for k, v := range pokers {
		vtype := types[v.Type()]
		vtype = vtype + 1
		types[v.Type()] = vtype
		if k == 0 {
			continue
		}
		if pokers[k].Value()-pokers[k-1].Value() == 1 {
			continue
		}
		return false
	}
	if pokers[len(pokers)-1].Value() == poker.V2 {
		return false
	}
	typenum := 0
	for _, v := range types {
		if v != 0 {
			typenum++
		}
	}
	if typenum == 1 {
		return false
	}
	return true
}

// NewSequence NewSequence
func NewSequence(pokers []poker.IPoker) Sequence {
	return Sequence{
		pokers: pokers,
	}
}

// Type Type
func (p Sequence) Type() int {
	return TypeSequence
}

// Value Value
func (p Sequence) Value() int {
	return p.pokers[0].Value()
}

// Pokers Pokers
func (p Sequence) Pokers() []poker.IPoker {
	return p.pokers
}

// IsBoom IsBoom
func (p Sequence) IsBoom() bool {
	return false
}

// Length Length
func (p Sequence) Length() int {
	return len(p.pokers)
}
