package pokertype

import (
	"ddz/game/poker"
)

func sortPokers(pokers []poker.IPoker, l int, r int) []poker.IPoker {
	return nil
}

func compare(p1 poker.IPoker, p2 poker.IPoker) int {
	if p1.Value() > p2.Value() {
		return 1
	} else if p1.Value() < p2.Value() {
		return -1
	}
	if p1.Type() > p2.Type() {
		return 1
	}
	if p1.Type() < p2.Type() {
		return -1
	}
	return 0
}
