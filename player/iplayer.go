package player

import (
	"ddz/message"
	"ddz/poker"
	"fmt"
)

// Status 状态
type Status = int8

const (
	Break      = -1 // 断线
	NotPrepare = 0  // 未准备
	Prepare    = 1  // 已准备
	Playing    = 2  // 进行中
)

// IPlayer 玩家
type IPlayer interface {
	Left() []poker.IPoker   // 剩余牌
	Played() []poker.IPoker // 出过得牌
	Status() Status         // 状态
	Sit(int)                // 坐在牌桌了
}

// Player 玩家
type Player struct {
	i            int                  // 座位号
	pokersLeft   []poker.IPoker       // 剩余牌
	pokersPlayed []poker.IPoker       // 出过的牌
	status       Status               // 状态
	recvChannel  chan message.Message // 接受频道
	sencChannel  chan message.Message // 发送频道
}

// NewPlayer NewPlayer
func NewPlayer(recvc chan message.Message, sendc chan message.Message) Player {
	return Player{
		pokersLeft:   []poker.IPoker{},
		pokersPlayed: []poker.IPoker{},
		status:       NotPrepare,
	}
}

// Left Left
func (p Player) Left() []poker.IPoker {
	return p.pokersLeft
}

// Played Played
func (p Player) Played() []poker.IPoker {
	return p.pokersPlayed
}

// Status Status
func (p Player) Status() Status {
	return p.status
}

// Sit 坐在牌桌了
func (p Player) Sit(i int) {
	go func() {
		select {
		case msg := <-p.recvChannel:
			fmt.Println(msg)
		}
	}()
}

// Play 出牌
func (p Player) Play(pokers []poker.IPoker) {
	p.sencChannel <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerPlay,
		Chat:          "",
		PlayerCurrent: p.i,
		Pokers:        pokers,
	}
}
