package player

import (
	"ddz/message"
	"ddz/poker"
	"errors"
	"fmt"
	"strconv"
)

// Status 状态
type Status = int8

const (
	Break   = -1 // 断线
	Sit     = 0  // 已就坐
	Prepare = 1  // 已准备
	Playing = 2  // 进行中
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
		status:       Sit,
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
	// player设置座位号
	p.i = i

	// 通知牌桌已就坐
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerSit,
		Chat:          "",
		PlayerCurrent: i,
		PlayerTurn:    0,
		Pokers:        nil,
	}
}

// Play 出牌
func (p *Player) Play(pokers []poker.IPoker) error {
	for _, pk := range pokers {
		if !p.pokerIsMine(pk) {
			// TODO debug info
			fmt.Println("system: " + p.name + "没有这张牌 - " + pk.Show())
			return errors.New("没有这张牌")
		}
	}
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerPlay,
		Chat:          "",
		PlayerCurrent: p.i,
		Pokers:        pokers,
	}
	return nil
}

// PlayNone 不出
func (p *Player) PlayNone() error {
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerPlay,
		Chat:          "",
		PlayerCurrent: p.i,
		Pokers:        []poker.IPoker{},
	}
	return nil
}

// PlayAll 全部出掉
func (p *Player) PlayAll() error {
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerPlay,
		Chat:          "",
		PlayerCurrent: p.i,
		Pokers:        p.pokersLeft,
	}
	return nil
}

func (p *Player) pokerIsMine(pk poker.IPoker) bool {
	for _, pl := range p.pokersLeft {
		if pl.Type() == pk.Type() && pl.Value() == pk.Value() {
			return true
		}
	}
	return false
}

// SetRevc 玩家设置接受牌桌信息的管道
// 设置完立刻监听
func (p *Player) SetRevc(recv chan message.Message) {
	p.recv = recv

	go func() {
		for {
			select {
			case msg := <-p.recv:
				switch msg.T {
				// 聊天信息
				case message.TypeChat:
					fmt.Println(p.name + " recive: [" + msg.Chat + "]")
				case message.TypeNotice:
					fmt.Println(p.name + " recive: [" + msg.Chat + "]")
				// 游戏中
				case message.TypeRuler:
					switch msg.ST {
					// 已就坐
					case message.SubTypeRulerSit:
						fmt.Println(p.name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家已就坐]")
					// 已准备
					case message.SubTypeRulerReady:
						fmt.Println(p.name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家已准备]")
					// 洗牌中
					case message.SubTypeRulerShuffle:
						fmt.Println(p.name + " recive: [洗牌中]")
					// 发牌
					case message.SubTypeRulerReal:
						showpokers := ""
						for _, v := range msg.Pokers {
							showpokers += v.Show()
						}
						fmt.Println(p.name + " recive: [发牌:( " + showpokers + " )]")
						p.pokersLeft = msg.Pokers
					// 出牌
					case message.SubTypeRulerPlay:
						showpokers := ""
						for _, v := range msg.Pokers {
							showpokers += v.Show()
						}
						fmt.Println(p.name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家出牌:( " + showpokers + " )]")
					case message.SubTypeRulerChangePlayer:
						fmt.Println(p.name + " recive: [现在轮到" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家出牌]")
					case message.SubTypeRulerEnd:
						fmt.Println(p.name + " recive: [本局游戏结束]")
						// TODO 一些玩家本局数据恢复的设置
					}
				}
			}
		}
	}()
}

// SetSend 玩家设置发送牌桌信息的管道
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

// Prepare 准备
func (p *Player) Prepare() error {
	if p.status != Sit {
		return errors.New("状态出错，不能准备开始")
	}

	// 玩家状态改为已准备
	p.status = Prepare

	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerReady,
		Chat:          "",
		PlayerCurrent: p.i,
		Pokers:        nil,
	}
	return nil
}
