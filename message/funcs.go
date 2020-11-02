package message

import (
	"encoding/json"
	"fmt"
)

// Decode 解码
func Decode(bts []byte) (Message, error) {
	msg := Message{}
	if err := json.Unmarshal(bts, &msg); err != nil {
		fmt.Println(err.Error())
		msg.T = TypeNone
		return msg, nil
	}
	return msg, nil
}

// Encode 编码
func Encode(msg Message) []byte {
	bts, _ := json.Marshal(msg)
	return bts
}
