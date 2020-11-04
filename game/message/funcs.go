package message

import (
	"ddz/game/poker"
	"encoding/json"
	"fmt"
)

// Decode 解码
func Decode(bts []byte) (Message, error) {
	var _msg = struct {
		T                  Type          `json:"type"`
		ST                 SubType       `json:"sub_type"`
		Chat               string        `json:"chat"`
		PlayerCurrent      int           `json:"player_current"`
		PlayerTurn         int           `json:"player_turn"`
		Pokers             []poker.Poker `json:"pokers"`
		TableIndex         int           `json:"table_index"`
		TablePositionIndex int           `json:"table_position_index"`
		NewCircle          bool          `json:"new_circle"`
		Data               interface{}   `json:"data"`
	}{}

	msg := Message{}
	if err := json.Unmarshal(bts, &_msg); err != nil {
		fmt.Println(err.Error())
		msg.T = TypeNone
		return msg, nil
	}
	msg.T = _msg.T
	msg.ST = _msg.ST
	msg.Chat = _msg.Chat
	msg.PlayerCurrent = _msg.PlayerCurrent
	msg.PlayerTurn = _msg.PlayerTurn
	var pks = []poker.IPoker{}
	for _, v := range _msg.Pokers {
		pks = append(pks, v)
	}
	msg.Pokers = pks
	msg.TableIndex = _msg.TableIndex
	msg.TablePositionIndex = _msg.TablePositionIndex
	msg.NewCircle = _msg.NewCircle
	msg.Data = _msg.Data
	return msg, nil
}

// Encode 编码
func Encode(msg Message) []byte {
	bts, _ := json.Marshal(msg)
	return bts
}
