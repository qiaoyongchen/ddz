package poker

import (
	"errors"
	"math/rand"
	"time"
)

// SubPokers 从 origin 中剪掉 toSub 中的牌
// 如果某张牌在 toSub 中有,但是在 origin 中没有则报错
func SubPokers(origin []IPoker, toSub []IPoker) ([]IPoker, error) {
	left := make([]IPoker, 0, len(origin)-len(toSub))
	for _, pk := range origin {
		if len(toSub) == 0 {
			left = append(left, pk)
			continue
		}
		inFlag := false
		for _, pkk := range toSub {
			if pk.CompareTo(pkk) == 0 {
				inFlag = true
				break
			}
		}
		if !inFlag {
			left = append(left, pk)
		}
	}
	if (len(left) + len(toSub)) != len(origin) {
		return nil, errors.New("计算牌出错")
	}
	return left, nil
}

// Contain origin 是否包含 toSub
func Contain(origin []IPoker, toSub []IPoker) bool {
	_, err := SubPokers(origin, toSub)
	return err != nil
}

// Shuffle 打乱一组牌
func Shuffle(pokers []IPoker) {
	max := len(pokers) - 1
	for i := 0; i <= 1000; i++ {
		rand.Seed(time.Now().UnixNano())
		fst := rand.Intn(max + 1)
		rand.Seed(time.Now().UnixNano() + int64(fst))
		snd := rand.Intn(max + 1)
		pokers[fst], pokers[snd] = pokers[snd], pokers[fst]
	}
}

// ShowPokers 字符形式显示牌
func ShowPokers(pokers []IPoker) string {
	showpokers := ""
	for _, v := range pokers {
		showpokers += v.Show()
	}
	return showpokers
}
