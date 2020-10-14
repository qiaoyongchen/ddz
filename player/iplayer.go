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
	Left() []poker.IPoker         // 剩余牌
	Played() []poker.IPoker       // 出过得牌
	Status() Status               // 状态
	Sit(int)                      // 坐在牌桌了
	SetRevc(chan message.Message) // 设置接受管道
	SetSend(chan message.Message) // 设置接受管道
	Name() string                 // 玩家名字
}

// Player 玩家
type Player struct {
	name         string               // 玩家姓名
	i            int                  // 座位号
	pokersLeft   []poker.IPoker       // 剩余牌
	pokersPlayed []poker.IPoker       // 出过的牌
	status       Status               // 状态
	recv         chan message.Message // 接受频道
	send         chan message.Message // 发送频道
}

// NewPlayer NewPlayer
func NewPlayer(name string) *Player {
	return &Player{
		pokersLeft:   []poker.IPoker{},
		pokersPlayed: []poker.IPoker{},
		status:       NotPrepare,
		name:         name,
	}
}

// Left Left
func (p *Player) Left() []poker.IPoker {
	return p.pokersLeft
}

// Played Played
func (p *Player) Played() []poker.IPoker {
	return p.pokersPlayed
}

// Status Status
func (p *Player) Status() Status {
	return p.status
}

// Sit 坐在牌桌了
func (p *Player) Sit(i int) {
	go func() {
		select {
		case msg := <-p.recv:
			fmt.Println(p.name + " recive: [" + msg.Chat + "]")
		}
	}()
}

// Play 出牌
func (p *Player) Play(pokers []poker.IPoker) {
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerPlay,
		Chat:          "",
		PlayerCurrent: p.i,
		Pokers:        pokers,
	}
}

// SetRevc SetRevc
func (p *Player) SetRevc(recv chan message.Message) {
	p.recv = recv
}

// SetSend SetSend
func (p *Player) SetSend(send chan message.Message) {
	p.send = send
}

// Name Name
func (p *Player) Name() string {
	return p.name
}

// Chat Chat
func (p *Player) Chat(content string) {
	p.send <- message.Message{
		T:             message.TypeChat,
		ST:            0,
		Chat:          p.name + " say: " + content,
		PlayerCurrent: 0,
		PlayerTurn:    0,
		Pokers:        nil,
	}
}
