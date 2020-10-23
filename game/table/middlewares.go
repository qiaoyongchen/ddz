package table

import "ddz/message"

// 出牌时检查是否轮到该用户出牌
func forPlayIsMyTurn(p processorFunc) processorFunc {
	return func(t *Table, msg message.Message) {
		if message.SubTypeRulerPlay == msg.ST && t.playerCurrent != msg.PlayerCurrent {
			t.sendone(msg.PlayerCurrent, message.Message{
				T:             message.TypeNotice,
				ST:            0,
				Chat:          "system: 现在没有轮到" + t.players[msg.PlayerCurrent].Name() + "出牌",
				PlayerCurrent: msg.PlayerCurrent,
				PlayerTurn:    msg.PlayerCurrent,
				Pokers:        msg.Pokers,
			})
			return
		}
		p(t, msg)
	}
}

// 出牌时检查是否比上一轮的牌大
// 如果上一轮最大牌是自己说明出了一圈没人比该玩家大，重新轮到该玩家出牌了，不需要检查规则
// 如果上一轮最大牌不是自己说明有人接牌了，需要检查出牌规则一致
func forPlayBigThanLast(p processorFunc) processorFunc {
	return func(t *Table, msg message.Message) {
		// 要不起
		if len(msg.Pokers) == 0 {
			t.broadcast(msg)
			return
		}

		// TODO

		p(t, msg)
	}
}
