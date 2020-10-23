package ruler

import (
	"ddz/game/poker"
)

// IRuler 规则
type IRuler interface {
	Check([]poker.IPoker, []poker.IPoker) (bool, error)
}

// Ruler 规则
type Ruler struct {
}

// Check 检查这一把牌和上一把牌并返回大小(now > last),如果出错就报错
func (p Ruler) Check(now []poker.IPoker, last []poker.IPoker) (bool, error) {
	return true, nil
}

// NewRuler NewRuler
func NewRuler() Ruler {
	return Ruler{}
}
