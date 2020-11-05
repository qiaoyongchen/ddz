package message

import (
	"ddz/game/poker"
	"strconv"

	"github.com/gorilla/websocket"
)

// Type 消息类型
type Type = int8

const (
	// TypeNone 错误的消息
	TypeNone = -1
	// TypeNotice 通知消息类型(比如断线了进行广播,牌出错了进行广播)
	TypeNotice = 0
	// TypeChat 聊天消息类型(用于玩家聊天)
	TypeChat = 1
	// TypeRuler 游戏规则类型(用于游戏规则,比如:已就坐，准备，洗牌，发牌，出牌，结束)
	TypeRuler = 2
	// TypeRoom 在房间里使用的类型
	TypeRoom = 3
)

// SubType 消息子类型
type SubType = int8

const (
	// SubTypeNoticeBreak 断线了
	SubTypeNoticeBreak Type = 0
	// SubTypeNoticeRelink 断线又重连了
	SubTypeNoticeRelink Type = 1
	// SubTypeNoticeError 报错消息
	SubTypeNoticeError Type = 2 // 报错消息
	// SubTypeRulerSit 已就坐
	SubTypeRulerSit SubType = 0
	// SubTypeRulerReady 用户已准备
	SubTypeRulerReady SubType = 1
	// SubTypeRulerShuffle 正在洗牌
	SubTypeRulerShuffle SubType = 2
	// SubTypeRulerReal 正在发牌
	SubTypeRulerReal SubType = 3
	// SubTypeRulerPlay 出牌
	SubTypeRulerPlay SubType = 4
	// SubTypeRulerEnd 游戏结束
	SubTypeRulerEnd SubType = 5
	// SubTypeRulerChangePlayer 改变出牌人
	SubTypeRulerChangePlayer SubType = 6
	// SubTypeRulerWinner 获胜
	SubTypeRulerWinner SubType = 7
	// SubTypeRoomInfo 显示房间信息
	SubTypeRoomInfo SubType = 0
	// SubTypeGetRoomInfo 请求获取显示房间信息
	SubTypeGetRoomInfo SubType = 1
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
	NewCircle          bool           `json:"new_circle"`           // 新的一轮出牌
	Data               interface{}    `json:"data"`                 // 携带的其他数据
}

// Send 发送消息
func (p Message) Send(ch chan<- Message) {
	ch <- p
}

// WriteConn 写入到websocket
func (p Message) WriteConn(conn *websocket.Conn) error {
	return conn.WriteMessage(websocket.TextMessage, Encode(p))
}

// String 消息转换为字符形式(用于测试)
func (p Message) String() string {
	switch p.T {
	case TypeChat:
		return p.Chat
	case TypeNotice:
		return p.Chat
	case TypeRuler:
		switch p.ST {
		case SubTypeRulerSit:
			return strconv.Itoa(p.PlayerCurrent) + "号位置玩家已就坐"
		case SubTypeRulerReady:
			return strconv.Itoa(p.PlayerCurrent) + "号位置玩家已准备]"
		case SubTypeRulerShuffle:
			return "洗牌中"
		case SubTypeRulerReal:
			return "发牌:( " + poker.ShowPokers(p.Pokers) + " )"
		case SubTypeRulerPlay:
			if len(p.Pokers) == 0 {
				return "不要"
			}
			return strconv.Itoa(p.PlayerCurrent) + "号位置玩家出牌:( " + poker.ShowPokers(p.Pokers) + " )"
		case SubTypeRulerChangePlayer:
			return "现在轮到" + strconv.Itoa(p.PlayerCurrent) + "号位置玩家出牌"
		case SubTypeRulerWinner:
			return strconv.Itoa(p.PlayerCurrent) + "号位置玩家获胜"
		case SubTypeRulerEnd:
			return "本局游戏结束"
		}
	}
	return "unknow message"
}

// GenMessageWinner 获胜消息
func GenMessageWinner(winnerIndex int) Message {
	return Message{
		T:             TypeRuler,
		ST:            SubTypeRulerWinner,
		PlayerCurrent: winnerIndex,
	}
}

// GenMessageEnd 一局结束消息
func GenMessageEnd() Message {
	return Message{
		T:             TypeRuler,
		ST:            SubTypeRulerEnd,
		PlayerCurrent: -1,
	}
}

// GenMessageChangePlayer 更换出派人消息
func GenMessageChangePlayer(nextPlayerIndex int, newCircle bool) Message {
	return Message{
		T:             TypeRuler,
		ST:            SubTypeRulerChangePlayer,
		PlayerCurrent: nextPlayerIndex,
		NewCircle:     newCircle,
	}
}

// GenMessageReal 发牌消息
func GenMessageReal(playerIndex int, pokers []poker.IPoker) Message {
	return Message{
		T:             TypeRuler,
		ST:            SubTypeRulerReal,
		PlayerCurrent: playerIndex,
		Pokers:        pokers,
	}
}

// GenMessageShuffle 洗牌信息
func GenMessageShuffle() Message {
	return Message{
		T:  TypeRuler,
		ST: SubTypeRulerShuffle,
	}
}

// GenMessageNoticeError 通知消息
func GenMessageNoticeError(content string) Message {
	return Message{
		T:    TypeNotice,
		ST:   SubTypeNoticeError,
		Chat: content,
	}
}

// GenMessageRulerSit 玩家就坐消息
func GenMessageRulerSit(playerIndex int) Message {
	return Message{
		T:             TypeRuler,
		ST:            SubTypeRulerSit,
		PlayerCurrent: playerIndex,
	}
}

// GenMessageChat 聊天消息
func GenMessageChat(playerIndex int, content string) Message {
	return Message{
		T:             TypeChat,
		Chat:          content,
		PlayerCurrent: playerIndex,
	}
}

// GenMessageRoomInfo 房间信息
func GenMessageRoomInfo(data interface{}) Message {
	return Message{
		T:    TypeRoom,
		ST:   SubTypeRoomInfo,
		Data: data,
	}
}
