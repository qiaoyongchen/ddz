package table

import (
	"ddz/game/message"
	"ddz/game/poker"
	"ddz/game/proc1"
)

// 聊天
func proc4Chat(t *Table) proc1.ProcessorFunc {
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
		case message.SubTypeRulerPlay:
			t.broadcast(msg)
			if len(msg.Pokers) == len(t.players[msg.PlayerCurrent].Left()) {
				t.end(msg.PlayerCurrent)
				return
			}

			player := t.Players()[msg.PlayerCurrent]
			newleft, err := poker.SubPokers(player.Left(), msg.Pokers)
			if err != nil {
				t.sendone(msg.PlayerCurrent, message.GenMessageNoticeError("出的牌不对"))
				return
			}
			player.SetLeft(newleft)
			player.SetPlayed(append(player.Played(), msg.Pokers...))

			t.playerMaxPokers = msg.PlayerCurrent
			t.maxPokers = msg.Pokers
			t.nextPlayer(t.playerCurrent)
		}
	}
}
