package ruler

import (
	"ddz/game/poker"
	"fmt"
	"testing"
)

var testPokers = map[string][]poker.IPoker{
	"单张": []poker.IPoker{
		poker.NewPoker(poker.TypeSpade, poker.V3),
	},
	// 对子
	"对子": []poker.IPoker{
		poker.NewPoker(poker.TypeHeart, poker.V3),
		poker.NewPoker(poker.TypeSpade, poker.V3),
	},
	"三代一": []poker.IPoker{
		poker.NewPoker(poker.TypeSpade, poker.V3),
		poker.NewPoker(poker.TypeHeart, poker.V3),
		poker.NewPoker(poker.TypeClub, poker.V3),
		poker.NewPoker(poker.TypeClub, poker.V4),
	},
	"三代二": []poker.IPoker{
		poker.NewPoker(poker.TypeSpade, poker.V3),
		poker.NewPoker(poker.TypeHeart, poker.V3),
		poker.NewPoker(poker.TypeClub, poker.V3),
		poker.NewPoker(poker.TypeClub, poker.V4),
		poker.NewPoker(poker.TypeHeart, poker.V4),
	},
	"顺子": []poker.IPoker{
		poker.NewPoker(poker.TypeSpade, poker.V3),
		poker.NewPoker(poker.TypeClub, poker.V4),
		poker.NewPoker(poker.TypeClub, poker.V5),
		poker.NewPoker(poker.TypeClub, poker.V6),
		poker.NewPoker(poker.TypeClub, poker.V7),
	},
	"炸": []poker.IPoker{
		poker.NewPoker(poker.TypeSpade, poker.V3),
		poker.NewPoker(poker.TypeHeart, poker.V3),
		poker.NewPoker(poker.TypeClub, poker.V3),
		poker.NewPoker(poker.TypeDiamond, poker.V3),
	},
	"同花顺": []poker.IPoker{
		poker.NewPoker(poker.TypeSpade, poker.V3),
		poker.NewPoker(poker.TypeSpade, poker.V4),
		poker.NewPoker(poker.TypeSpade, poker.V5),
		poker.NewPoker(poker.TypeSpade, poker.V6),
		poker.NewPoker(poker.TypeSpade, poker.V7),
	},
	"王炸": []poker.IPoker{
		poker.NewPoker(poker.TypeLittleJoker, poker.TypeLittleJoker),
		poker.NewPoker(poker.TypeBigJoker, poker.VBigJoker),
	},
}

func TestPkt(t *testing.T) {
	r := NewRuler()

	for k, v := range testPokers {
		fmt.Println(k)

		r.Sort(v)
		fmt.Println(v)

		pkt, pkte := r.getPokerType(v)
		fmt.Println(pkt.Type(), pkt.Value(), "是否是炸弹:", pkt.IsBoom(), pkte)
		fmt.Println("---------------------------------")
		fmt.Println()
		fmt.Println()
	}
}

func TestCheck(t *testing.T) {
	r := NewRuler()
	rst, rste := r.Check([]poker.IPoker{poker.NewPoker(1, 3)}, []poker.IPoker{})
	fmt.Println("")
	fmt.Println(rst)
	fmt.Println(rste)
}
