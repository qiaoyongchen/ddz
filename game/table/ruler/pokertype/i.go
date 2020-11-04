package pokertype

import (
	"ddz/game/poker"
)

// Type 牌型的值
type Type = int

const (
	// TypeJokerBoom 王炸
	TypeJokerBoom Type = 1000
	// TypeFlushBoom 同花顺
	TypeFlushBoom Type = 100
	// TypeFourBoom 四个炸
	TypeFourBoom Type = 10
	// TypeThreeWithTwo 三带二
	TypeThreeWithTwo Type = 5
	// TypeThreeWithOne 三带一
	TypeThreeWithOne Type = 4
	// TypeSequence 顺子
	TypeSequence Type = 3
	// TypePair 对
	TypePair Type = 2
	// TypeSingle 单张
	TypeSingle Type = 1
	// TypeEmpty 空
	TypeEmpty Type = 0
)

// PokerType 牌型
type PokerType interface {
	Type() int
	Value() int
	Pokers() []poker.IPoker
	IsBoom() bool
	Length() int
}
