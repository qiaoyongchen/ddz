package poker

import "strconv"

type Type = uint

const (
	TypeLittleJoker = 10000 // 小王
	TypeBigJoker    = 10001 // 大王
	TypeSpade       = 3     // 黑桃
	TypeHeart       = 2     // 红桃
	TypeClub        = 1     // 梅花
	TypeDiamond     = 0     // 方块
)

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
}

type Poker struct {
	t Type
	v Value
}

func NewPoker(t Type, v Value) Poker {
	return Poker{
		t: t,
		v: v,
	}
}

func (p Poker) Type() Type {
	return p.t
}

func (p Poker) Value() Value {
	return p.v
}

func (p Poker) Show() string {
	showstring := ""
	if p.t == TypeSpade {
		showstring += "[♠"
	} else if p.t == TypeHeart {
		showstring += "[♥"
	} else if p.t == TypeClub {
		showstring += "[♣"
	} else if p.t == TypeDiamond {
		showstring += "[♦"
	} else if p.t == TypeLittleJoker {
		showstring += "[Little🃏]"
		return showstring
	} else if p.t == TypeBigJoker {
		showstring += "[Big🃏]"
		return showstring
	}

	if p.v >= 3 && p.v <= 10 {
		return "-" + showstring + strconv.Itoa(p.v) + "]"
	} else if p.v == 11 {
		return " " + showstring + "-J" + "]"
	} else if p.v == 12 {
		return " " + showstring + "-Q" + "]"
	} else if p.v == 13 {
		return " " + showstring + "-K" + "]"
	} else if p.v == 14 {
		return " " + showstring + "-A" + "]"
	} else if p.v == 15 {
		return " " + showstring + "-2" + "]"
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
