package pokertype

import "ddz/game/poker"

// FlushBoom FlushBoom
type FlushBoom struct {
	pokers []poker.IPoker
}

// IsFlushBoom 是否是同花顺
func IsFlushBoom(pokers []poker.IPoker) bool {
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
	if typenum != 1 {
		return false
	}
	return true
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

// Length Length
func (p FlushBoom) Length() int {
	return len(p.pokers)
}
