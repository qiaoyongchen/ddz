package ruler

import (
	"ddz/game/poker"
	"ddz/game/table/ruler/pokertype"
	"errors"
)

// IRuler 规则
type IRuler interface {
	Check([]poker.IPoker, []poker.IPoker) (int, error)
	Shuffle([]poker.IPoker)
}

// Ruler 规则
type Ruler struct {
}

// NewRuler NewRuler
func NewRuler() Ruler {
	return Ruler{}
}

// Check 检查这一把牌和上一把牌并返回大小(now > last),如果出错就报错, 0 相等, 1 大于, -1 小于
func (p Ruler) Check(p1 []poker.IPoker, p2 []poker.IPoker) (int, error) {
	p.Sort(p1)
	p.Sort(p2)
	p1t, p1terr := p.getPokerType(p1)
	if p1terr != nil {
		return 0, p1terr
	}
	p2t, p2terr := p.getPokerType(p1)
	if p2terr != nil {
		return 0, p2terr
	}
	if !p1t.IsBoom() && !p2t.IsBoom() {
		if p1t.Type() != p2t.Type() {
			return 0, errors.New("牌型不匹配")
		}
		if p1t.Value() > p2t.Value() {
			return 1, nil
		} else if p1t.Value() == p2t.Value() {
			return 0, nil
		} else {
			return -1, nil
		}
	} else if !p1t.IsBoom() && p2t.IsBoom() {
		return -1, nil
	} else if p1t.IsBoom() && !p2t.IsBoom() {
		return 1, nil
	} else if p1t.IsBoom() && p2t.IsBoom() {
		if p1t.Type() > p2t.Type() {
			return 1, nil
		} else if p1t.Type() < p2t.Type() {
			return -1, nil
		} else {
			// 炸弹的类型不一样时比较炸弹级别
			if p1t.Type() != pokertype.TypeFlushBoom {
				if p1t.Value() > p2t.Value() {
					return 1, nil
				} else if p1t.Value() < p2t.Value() {
					return -1, nil
				} else {
					return 0, nil
				}
			}
			// 炸弹的类型一样时比较炸弹值,但同花顺时特殊
			if p1t.Type() == pokertype.TypeFlushBoom {
				if p1t.Length() != p2t.Length() {
					return 0, errors.New("牌型不匹配")
				}
				if p1t.Value() > p2t.Value() {
					return 1, nil
				} else if p1t.Value() < p2t.Value() {
					return -1, nil
				} else {
					return 0, nil
				}
			}
		}
	}
	return 0, errors.New("牌型比对错误")
}

// Shuffle 洗牌
func (p Ruler) Shuffle(pokers []poker.IPoker) {
	poker.Shuffle(pokers)
}

func (p Ruler) getPokerType(ps []poker.IPoker) (pokertype.PokerType, error) {
	if pokertype.IsSingle(ps) {
		return pokertype.NewSingle(ps), nil
	}
	if pokertype.IsPair(ps) {
		return pokertype.NewPair(ps), nil
	}
	if pokertype.IsThreeWithOne(ps) {
		return pokertype.NewThreeWithOne(ps), nil
	}
	if pokertype.IsThreeWithTwo(ps) {
		return pokertype.NewThreeWithTwo(ps), nil
	}
	if pokertype.IsSequence(ps) {
		return pokertype.NewSequence(ps), nil
	}
	if pokertype.IsFourBoom(ps) {
		return pokertype.NewFourBoom(ps), nil
	}
	if pokertype.IsFlushBoom(ps) {
		return pokertype.NewFlushBoom(ps), nil
	}
	if pokertype.IsJokerBoom(ps) {
		return pokertype.NewJokerBoom(ps), nil
	}
	return nil, errors.New("出牌不符合规格")
}

// Sort 排序扑克牌
func (p Ruler) Sort(ps []poker.IPoker) {
	sortPokers(ps, 0, len(ps)-1)
}
