package table

import (
	"ddz/message"
	"time"
)

func chatProcessor(t *Table, msg message.Message) {
	t.broadcast(msg)
}

func rulerProcessor(t *Table, msg message.Message) {
	switch msg.ST {
	// 已就坐
	case message.SubTypeRulerSit:
		t.broadcast(msg)
	// 已准备
	case message.SubTypeRulerReady:
		t.broadcast(msg)
		// 每次收到玩家已准备信息都检查一次是否可以开始打牌了
		if t.allReady() {
			// 洗牌
			t.shuffle()
			// 停顿1秒，给客户端停留，显示洗牌画面
			time.Sleep(time.Second)
			// 发牌
			t.real()
			// 指定第一个出牌玩家
			t.nextPlayer(-1)
		}
	case message.SubTypeRulerPlay:
		// 检查现在是否轮到该用户出牌
		if t.playerCurrent != msg.PlayerCurrent {
			showpokers := ""
			for _, v := range msg.Pokers {
				showpokers += v.Show()
			}
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

		// 出牌广播
		t.broadcast(msg)

		// 每次玩家出完牌都检查是不是游戏结束了
		// 如果结束就广播结束通知，否则切换到下一个玩家出牌
		if len(msg.Pokers) == len(t.players[msg.PlayerCurrent].Left()) {
			t.end()
		} else {
			t.nextPlayer(t.playerCurrent)
		}
	}
}
