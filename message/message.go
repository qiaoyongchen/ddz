package message

import (
	"ddz/game/poker"
	"encoding/json"
)

type Type = int8

const (
	TypeNone   = -1 // 错误的消息
	TypeNotice = 0  // 通知消息类型(比如断线了进行广播,牌出错了进行广播)
	TypeChat   = 1  // 聊天消息类型(用于玩家聊天)
	TypeRuler  = 2  // 游戏规则类型(用于游戏规则,比如:已就坐，准备，洗牌，发牌，出牌，结束)
	TypeRoom   = 3  // 在房间里使用的类型
)

type SubType = int8

const (
	SubTypeNoticeBreak       Type    = 0 // 断线了
	SubTypeNoticeRelink      Type    = 1 // 断线又重连了
	SubTypeNoticeError       Type    = 2 // 报错消息
	SubTypeRulerSit          SubType = 0 // 已就坐
	SubTypeRulerReady        SubType = 1 // 用户已准备
	SubTypeRulerShuffle      SubType = 2 // 正在洗牌
	SubTypeRulerReal         SubType = 3 // 正在发牌
	SubTypeRulerPlay         SubType = 4 // 出牌
	SubTypeRulerEnd          SubType = 5 // 游戏结束
	SubTypeRulerChangePlayer SubType = 6 // 改变出牌人
	SubTypeRulerWinner       SubType = 7 // 获胜
	SubTypeRoomInfo          SubType = 0 // 显示房间信息
)

// Message 消息
type Message struct {
	T                  Type           `json:"type"`                 // type
	ST                 SubType        `json:"sub_type"`             // sub_type
	Chat               string         `json:"chat"`                 // chat 聊天内容
	PlayerCurrent      int            `json:"player_current"`       // 当前用户
	PlayerTurn         int            `json:"player_turn"`          // 轮到哪个用户出牌
	Pokers             []poker.IPoker `json:"pokers"`               // 牌
	TableIndex         int            `json:"table_index"`          // 第几桌
	TablePositionIndex int            `json:"table_position_index"` // 这桌第几个位置
	Data               interface{}    `json:"data"`                 // 携带的其他数据
}

// Decode 解码
func Decode(bts []byte) (Message, error) {
	msg := Message{}
	if err := json.Unmarshal(bts, &msg); err != nil {
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
