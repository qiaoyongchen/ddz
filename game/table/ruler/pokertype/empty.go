package pokertype

import "ddz/game/poker"

// Empty 空牌
type Empty struct {
}

// IsEmpty IsEmpty
func IsEmpty(pokers []poker.IPoker) bool {
	return len(pokers) == 0
}

// NewEmpty NewEmpty
func NewEmpty() Empty {
	return Empty{}
}

// Type Type
func (p Empty) Type() int {
	return TypeEmpty
}

// Value Value
func (p Empty) Value() int {
	return 0
}

// Pokers Pokers
func (p Empty) Pokers() []poker.IPoker {
	return []poker.IPoker{}
}

// IsBoom IsBoom
func (p Empty) IsBoom() bool {
	return false
}

// Length Length
func (p Empty) Length() int {
	return 0
}
