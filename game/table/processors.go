package table

import (
	"ddz/game/message"
	"ddz/game/poker"
	"ddz/game/proc1"
	"fmt"
)

// 聊天
func proc4Chat(t *Table) proc1.ProcessorFunc {
	return func(msg message.Message) {
		t.broadcast(msg)
	}
}

// 通知
func proc4Notice(t *Table) proc1.ProcessorFunc {
	return func(msg message.Message) {
		t.broadcast(msg)
	}
}

// 按规则玩牌中
func proc4Ruler(t *Table) proc1.ProcessorFunc {
	return func(msg message.Message) {
		switch msg.ST {

		case message.SubTypeRulerSit:
			t.broadcast(msg)

		case message.SubTypeRulerReady:
			t.players[msg.PlayerCurrent].Ready()
			t.broadcast(msg)
			t.ready()

		// 出牌
		case message.SubTypeRulerPlay:
			if len(msg.Pokers) == 0 {
				t.broadcast(msg)
				t.nextPlayer(msg.PlayerCurrent)
				return
			}

			player := t.Players()[msg.PlayerCurrent]
			newleft, err := poker.SubPokers(player.Left(), msg.Pokers)
			if err != nil {
				t.sendone(msg.PlayerCurrent, message.GenMessageNoticeError("出的牌不对: 出的牌不是你手中有的牌"))
				fmt.Println("出的牌不对: 出的牌不是你手中有的牌")
				fmt.Println("player left:", player.Left())
				fmt.Println("player play:", msg.Pokers)
				return
			}

			// 出完之后没有牌了,本局结束
			if len(newleft) == 0 {
				t.end(msg.PlayerCurrent)
				return
			}

			player.SetLeft(newleft)
			player.SetPlayed(append(player.Played(), msg.Pokers...))

			t.playerMaxPokers = msg.PlayerCurrent
			t.maxPokers = msg.Pokers

			// 如果是新的一轮，出玩牌就不是了
			if t.ifNewCircle {
				t.ifNewCircle = false
			}

			// 调试信息 ...
			fmt.Println("玩家", msg.PlayerCurrent, "剩余牌:", player.Left())
			fmt.Println("玩家", msg.PlayerCurrent, "出牌牌:", msg.Pokers)
			fmt.Println("当前最大牌:", t.maxPokers)
			fmt.Println("当前最大牌玩家", t.playerMaxPokers)

			t.broadcast(msg)
			t.nextPlayer(t.playerCurrent)
		}
	}
}
