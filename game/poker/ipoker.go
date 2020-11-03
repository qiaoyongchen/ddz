package poker

import "strconv"

// Type 类型
type Type = int

const (
	// TypeLittleJoker 小王
	TypeLittleJoker Type = 10000
	// TypeBigJoker 大王
	TypeBigJoker Type = 10001
	// TypeSpade 黑桃
	TypeSpade Type = 3
	// TypeHeart 红桃
	TypeHeart Type = 2
	// TypeClub 梅花
	TypeClub Type = 1
	// TypeDiamond 方块
	TypeDiamond Type = 0
)

// Value 值
type Value = int

const (
	V3           = 3     // 3
	V4           = 4     // 4
	V5           = 5     // 5
	V6           = 6     // 6
	V7           = 7     // 7
	V8           = 8     // 8
	V9           = 9     // 9
	V10          = 10    // 10
	VJ           = 11    // J
	VQ           = 12    // Q
	VK           = 13    // K
	VA           = 14    // A
	V2           = 15    // 2
	VLittleJoker = 10000 // 小王
	VBigJoker    = 10001 // 大王
)

type IPoker interface {
	Type() Type
	Value() Value
	Show() string
	CompareTo(IPoker) int
}

type Poker struct {
	T Type  `json:"type"`
	V Value `json:"value"`
}

func NewPoker(t Type, v Value) Poker {
	return Poker{
		T: t,
		V: v,
	}
}

func (p Poker) Type() Type {
	return p.T
}

func (p Poker) Value() Value {
	return p.V
}

func (p Poker) CompareTo(p2 IPoker) int {
	if p.Value() > p2.Value() {
		return 1
	} else if p.Value() < p2.Value() {
		return -1
	}
	if p.Type() > p2.Type() {
		return 1
	}
	if p.Type() < p2.Type() {
		return -1
	}
	return 0
}

func (p Poker) Show() string {
	showstring := ""
	if p.T == TypeSpade {
		showstring += "[♠"
	} else if p.T == TypeHeart {
		showstring += "[♥"
	} else if p.T == TypeClub {
		showstring += "[♣"
	} else if p.T == TypeDiamond {
		showstring += "[♦"
	} else if p.T == TypeLittleJoker {
		showstring += "[Little🃏]"
		return showstring
	} else if p.T == TypeBigJoker {
		showstring += "[Big🃏]"
		return showstring
	}

	if p.V >= 3 && p.V <= 10 {
		return showstring + "-" + strconv.Itoa(p.V) + "]"
	} else if p.V == 11 {
		return showstring + "-J" + "]"
	} else if p.V == 12 {
		return showstring + "-Q" + "]"
	} else if p.V == 13 {
		return showstring + "-K" + "]"
	} else if p.V == 14 {
		return showstring + "-A" + "]"
	} else if p.V == 15 {
		return showstring + "-2" + "]"
	}
	return "[unknown]"
}

// OnePack 一副牌
func OnePack() []IPoker {
	return []IPoker{
		NewPoker(TypeSpade, VA), NewPoker(TypeHeart, VA), NewPoker(TypeClub, VA), NewPoker(TypeDiamond, VA),
		NewPoker(TypeSpade, V2), NewPoker(TypeHeart, V2), NewPoker(TypeClub, V2), NewPoker(TypeDiamond, V2),
		NewPoker(TypeSpade, V3), NewPoker(TypeHeart, V3), NewPoker(TypeClub, V3), NewPoker(TypeDiamond, V3),
		NewPoker(TypeSpade, V4), NewPoker(TypeHeart, V4), NewPoker(TypeClub, V4), NewPoker(TypeDiamond, V4),
		NewPoker(TypeSpade, V5), NewPoker(TypeHeart, V5), NewPoker(TypeClub, V5), NewPoker(TypeDiamond, V5),
		NewPoker(TypeSpade, V6), NewPoker(TypeHeart, V6), NewPoker(TypeClub, V6), NewPoker(TypeDiamond, V6),
		NewPoker(TypeSpade, V7), NewPoker(TypeHeart, V7), NewPoker(TypeClub, V7), NewPoker(TypeDiamond, V7),
		NewPoker(TypeSpade, V8), NewPoker(TypeHeart, V8), NewPoker(TypeClub, V8), NewPoker(TypeDiamond, V8),
		NewPoker(TypeSpade, V9), NewPoker(TypeHeart, V9), NewPoker(TypeClub, V9), NewPoker(TypeDiamond, V9),
		NewPoker(TypeSpade, V10), NewPoker(TypeHeart, V10), NewPoker(TypeClub, V10), NewPoker(TypeDiamond, V10),
		NewPoker(TypeSpade, VJ), NewPoker(TypeHeart, VJ), NewPoker(TypeClub, VJ), NewPoker(TypeDiamond, VJ),
		NewPoker(TypeSpade, VQ), NewPoker(TypeHeart, VQ), NewPoker(TypeClub, VQ), NewPoker(TypeDiamond, VQ),
		NewPoker(TypeSpade, VK), NewPoker(TypeHeart, VK), NewPoker(TypeClub, VK), NewPoker(TypeDiamond, VK),
		NewPoker(TypeLittleJoker, VLittleJoker), NewPoker(TypeBigJoker, VBigJoker),
	}
}
