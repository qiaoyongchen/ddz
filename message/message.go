package message

import (
	"ddz/poker"
)

type Type = uint8

const (
	TypeNotice = 0 // 通知消息类型(比如断线了进行广播,牌出错了进行广播)
	TypeChat   = 1 // 聊天消息类型(用于玩家聊天)
	TypeRuler  = 2 // 游戏规则类型(用于游戏规则,比如:已就坐，准备，洗牌，发牌，出牌，结束)
)

type SubType = uint8

const (
	SubTypeNoticeBreak  = 0 // 断线了
	SubTypeNoticeRelink = 1 // 断线又重连了

	SubTypeRulerSit          = 0 // 已就坐
	SubTypeRulerReady        = 1 // 用户已准备
	SubTypeRulerShuffle      = 2 // 正在洗牌
	SubTypeRulerReal         = 3 // 正在发牌
	SubTypeRulerPlay         = 4 // 出牌
	SubTypeRulerEnd          = 5 // 游戏结束
	SubTypeRulerChangePlayer = 6 // 改变出牌人
	SubTypeRulerWinner       = 7 // 获胜
)

// Message 消息
type Message struct {
	T             Type           `json:"type"`           // type
	ST            SubType        `json:"sub_type"`       // sub_type
	Chat          string         `json:"chat"`           // chat 聊天内容
	PlayerCurrent int            `json:"player_current"` // 当前用户
	PlayerTurn    int            `json:"player_turn"`    // 轮到哪个用户出牌
	Pokers        []poker.IPoker `json:"pokers"`         // 牌
}
