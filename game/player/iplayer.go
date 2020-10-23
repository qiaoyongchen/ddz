package player

import (
	"ddz/game/poker"
	"ddz/message"
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
	SetSend(chan message.Message) // 设置发送管道
	Name() string                 // 玩家名字
}

// Player 玩家
type Player struct {
	Name         string               // 玩家姓名
	I            int                  // 座位号
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
		Name:         name,
	}
}

// Left Left
func (p *Player) Left() []poker.IPoker {
	return p.pokersLeft
}

// ShowLeft 打印剩余牌
func (p *Player) ShowLeft() {
	fmt.Println(showPokers(p.Left()))
}

func showPokers(pokers []poker.IPoker) string {
	showpokers := ""
	for _, v := range pokers {
		showpokers += v.Show()
	}
	return showpokers
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
	p.I = i
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
			return errors.New("没有这张牌")
		}
	}
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerPlay,
		Chat:          "",
		PlayerCurrent: p.I,
		Pokers:        pokers,
	}
	return nil
}

// PlayNone 不出
func (p *Player) PlayNone() error {
	return p.Play([]poker.IPoker{})
}

// PlayAll 全部出掉
func (p *Player) PlayAll() error {
	return p.Play(p.pokersLeft)
}

// PlayFirst 出第一张牌(测试)
func (p *Player) PlayFirst() error {
	return p.Play(p.pokersLeft[0:1])
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
				case message.TypeChat:
					fmt.Println(p.Name + " recive: [" + msg.Chat + "]")
				case message.TypeNotice:
					fmt.Println(p.Name + " recive: [" + msg.Chat + "]")
				case message.TypeRuler:
					switch msg.ST {
					case message.SubTypeRulerSit:
						fmt.Println(p.Name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家已就坐]")
					case message.SubTypeRulerReady:
						fmt.Println(p.Name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家已准备]")
					case message.SubTypeRulerShuffle:
						fmt.Println(p.Name + " recive: [洗牌中]")
					case message.SubTypeRulerReal:
						fmt.Println(p.Name + " recive: [发牌:( " + showPokers(msg.Pokers) + " )]")
						p.pokersLeft = msg.Pokers
					case message.SubTypeRulerPlay:
						if len(msg.Pokers) == 0 {
							fmt.Println(p.Name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家: 不要]")
							continue
						}
						newpkleft := p.pokersLeft
						for _, pk := range msg.Pokers {
							p.pokersPlayed = append(p.pokersPlayed, pk)
							for i, pkk := range newpkleft {
								if pkk.Type() == pk.Type() && pkk.Value() == pk.Value() {
									newpkleft = append(newpkleft[0:i], newpkleft[i+1:]...)
								}
							}
						}
						p.pokersLeft = newpkleft
						fmt.Println(p.Name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家出牌:( " + showPokers(msg.Pokers) + " )]")
					case message.SubTypeRulerChangePlayer:
						fmt.Println(p.Name + " recive: [现在轮到" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家出牌]")
					case message.SubTypeRulerWinner:
						fmt.Println(p.Name + " recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家获胜]")
					case message.SubTypeRulerEnd:
						fmt.Println(p.Name + " recive: [本局游戏结束]")
						p.pokersLeft = []poker.IPoker{}
						p.pokersPlayed = []poker.IPoker{}
						p.status = Sit
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

// GetName GetName
func (p *Player) GetName() string {
	return p.Name
}

// Chat Chat
func (p *Player) Chat(content string) {
	p.send <- message.Message{
		T:             message.TypeChat,
		ST:            0,
		Chat:          p.Name + " say: " + content,
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
	p.status = Prepare
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerReady,
		Chat:          "",
		PlayerCurrent: p.I,
		Pokers:        nil,
	}
	return nil
}
