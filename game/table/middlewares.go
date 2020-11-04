package table

import (
	"ddz/game/message"
	"ddz/game/proc1"
	"fmt"
	"strconv"
)

// 出牌时检查是否轮到该用户出牌
func mw4PlayIsMyTurn(t *Table) proc1.ProcessorMiddleware {
	return func(p proc1.Processor) proc1.Processor {
		return proc1.ProcessorFunc(func(msg message.Message) {
			if message.SubTypeRulerPlay == msg.ST && t.playerCurrent != msg.PlayerCurrent {
				t.sendone(msg.PlayerCurrent,
					message.GenMessageNoticeError("system: 现在没有轮到玩家"+strconv.Itoa(msg.PlayerCurrent)+"出牌"))
				return
			}
			p.Process(msg)
		})
	}
}

// 出牌时检查是否比上一轮的牌大
// 如果上一轮最大牌是自己说明出了一圈没人比该玩家大，重新轮到该玩家出牌了，不需要检查规则
// 如果上一轮最大牌不是自己说明有人接牌了，需要检查出牌规则一致
func mw4PlayBigThanLast(t *Table) proc1.ProcessorMiddleware {
	return func(p proc1.Processor) proc1.Processor {
		return proc1.ProcessorFunc(func(msg message.Message) {
			// 要不起
			var giveup = func() bool {
				return len(msg.Pokers) == 0
			}
			if t.ifNewCircle {
				// 新的一轮不能要不起, 随便出什么牌都可以
				// ps : 该玩家出的牌别人都不要，又轮到该玩家出牌了
				if giveup() {
					t.sendone(msg.PlayerCurrent, message.GenMessageNoticeError("新的一轮不能要不起, 随便出什么牌都可以"))
					return
				}
			}

			// 不是新的一轮

			// 1. 要么直接要不起
			if giveup() {
				p.Process(msg)
				return
			}

			// 2. 要么出牌要比本轮最大的牌要大
			rst, rstError := t.ruler.Check(msg.Pokers, t.maxPokers)
			if rstError != nil {
				t.sendone(msg.PlayerCurrent, message.GenMessageNoticeError("出牌错误: 牌型错误"))
				fmt.Println("出牌错误: 牌型错误 --> ", rstError.Error())
				return
			}

			// 出的牌没有大过本轮最大牌
			if rst <= 0 {
				t.sendone(msg.PlayerCurrent, message.GenMessageNoticeError("出牌错误: 出的牌太小啦"))
				fmt.Println("出牌错误: 出的牌太小啦")
				fmt.Println("max current:", t.maxPokers)
				fmt.Println("player current:", msg.Pokers)
				return
			}

			p.Process(msg)
		})
	}
}
