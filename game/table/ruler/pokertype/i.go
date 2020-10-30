package pokertype

import (
	"ddz/game/poker"
)

// Type 牌型的值
type Type = int

const (
	TypeJokerBoom    = 1000 // 王炸
	TypeFlushBoom    = 100  // 同花顺
	TypeFourBoom     = 10   // 四个炸
	TypeThreeWithTwo = 5    // 三带二
	TypeThreeWithOne = 4    // 三带一
	TypeSequence     = 3    // 顺子
	TypePair         = 2    // 对
	TypeSingle       = 1    // 单张
)

// PokerType 牌型
type PokerType interface {
	Type() int
	Value() int
	Pokers() []poker.IPoker
	IsBoom() bool
	Length() int
}
