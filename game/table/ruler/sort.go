package ruler

import (
	"ddz/game/poker"
)

// SortPokers SortPokers
func SortPokers(pokers []poker.IPoker, l int, r int) {
	ll := l
	rr := r
	if l >= r {
		return
	}
	flagPoker := pokers[l]
	trunLeft := false
	for l < r {
		if trunLeft {
			if compare(pokers[l], flagPoker) == 0 || compare(pokers[l], flagPoker) == 1 {
				pokers[r] = pokers[l]
				r--
				trunLeft = false
			} else {
				l++
				continue
			}
		} else {
			if compare(pokers[r], flagPoker) == -1 {
				pokers[l] = pokers[r]
				l++
				trunLeft = true
			} else {
				r--
				continue
			}
		}
	}
	pokers[l] = flagPoker
	SortPokers(pokers, ll, l)
	SortPokers(pokers, l+1, rr)
}

// 1 p1 > p2 , 0 p1 == p2, -1, p1 < p2
func compare(p1 poker.IPoker, p2 poker.IPoker) int {
	return p1.CompareTo(p2)
}
