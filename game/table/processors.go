package table

import (
	"ddz/message"
)

// 聊天
func chatProcessor(t *Table, msg message.Message) {
	t.broadcast(msg)
}

// 按规则玩牌中
func rulerProcessor(t *Table, msg message.Message) {
	switch msg.ST {
	// 已就坐
	case message.SubTypeRulerSit:
		t.broadcast(msg)
	// 已准备
	case message.SubTypeRulerReady:
		t.broadcast(msg)
		t.ready()
	// 出牌
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
