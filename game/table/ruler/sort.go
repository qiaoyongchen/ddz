package ruler

import (
	"ddz/game/poker"
)

// SortPokers SortPokers
func sortPokers(pokers []poker.IPoker, l int, r int) {
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
			if pokers[r].CompareTo(flagPoker) == -1 {
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

	sortPokers(pokers, ll, l)
	sortPokers(pokers, l+1, rr)
}
