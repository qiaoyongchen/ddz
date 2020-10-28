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
			if pokers[l].CompareTo(flagPoker) == 0 || pokers[l].CompareTo(flagPoker) == 1 {
				pokers[r] = pokers[l]
				r--
				trunLeft = false
			} else {
				l++
				continue
			}
		} else {
			if pokers[2].CompareTo(flagPoker) == -1 {
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
