package table

import (
	"ddz/game/proc1"
	"ddz/message"
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
			t.broadcast(msg)
			t.ready()
		case message.SubTypeRulerPlay:
			t.broadcast(msg)
			if len(msg.Pokers) == len(t.players[msg.PlayerCurrent].Left()) {
				t.end(msg.PlayerCurrent)
			} else {
				t.playerMaxPokers = msg.PlayerCurrent
				t.maxPokers = msg.Pokers
				t.nextPlayer(t.playerCurrent)
			}
		}
	}
}
