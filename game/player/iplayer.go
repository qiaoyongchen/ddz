package player

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"

	"ddz/game/message"
	"ddz/game/poker"
)

// Status 状态
type Status = int8

const (
	// Break 断线
	Break Status = -1
	// Sit 已就坐
	Sit Status = 0
	// Ready 已准备
	Ready Status = 1
	// Playing 进行中
	Playing Status = 2
)

// IPlayer 玩家
type IPlayer interface {
	Left() []poker.IPoker         // 剩余牌
	SetLeft([]poker.IPoker)       // 设置剩余牌数
	Played() []poker.IPoker       // 出过的牌
	SetPlayed([]poker.IPoker)     // 设置出过的牌
	Status() Status               // 状态
	Sit(int) error                // 坐在牌桌了
	Ready() error                 // 准备
	Restart()                     // 重新开始
	SetRevc(chan message.Message) // 设置接受管道
	SetSend(chan message.Message) // 设置发送管道
	GetName() string              // 玩家名字
}

// Player 玩家
type Player struct {
	Name         string               `json:"name"`  // 玩家姓名
	I            int                  `json:"index"` // 座位号
	pokersLeft   []poker.IPoker       ``             // 剩余牌
	pokersPlayed []poker.IPoker       ``             // 出过的牌
	status       Status               ``             // 状态
	recv         chan message.Message ``             // 接受频道
	send         chan message.Message ``             // 发送频道
	conn         *websocket.Conn      ``             // websocket 连接
}

// NewPlayer NewPlayer
func NewPlayer(name string, conn *websocket.Conn) *Player {
	return &Player{
		pokersLeft:   []poker.IPoker{},
		pokersPlayed: []poker.IPoker{},
		status:       Sit,
		Name:         name,
		conn:         conn,
	}
}

// Left 剩余牌
func (p *Player) Left() []poker.IPoker {
	return p.pokersLeft
}

// SetLeft 设置剩余牌
func (p *Player) SetLeft(pokers []poker.IPoker) {
	p.pokersLeft = pokers
}

// SetPlayed 设置出过的牌
func (p *Player) SetPlayed(pokers []poker.IPoker) {
	p.pokersPlayed = pokers
}

// Played 出过的牌
func (p *Player) Played() []poker.IPoker {
	return p.pokersPlayed
}

// Status 状态
func (p *Player) Status() Status {
	return p.status
}

// Sit 坐在牌桌
func (p *Player) Sit(i int) error {
	p.I = i
	p.status = Sit
	go p.startListening()
	message.GenMessageRulerSit(i).Send(p.send)
	return nil
}

// Ready 准备
func (p *Player) Ready() error {
	if p.status != Sit {
		return errors.New("状态出错，不能准备开始")
	}
	p.status = Ready
	return nil
}

// Restart 一轮结束，重新开始
func (p *Player) Restart() {
	p.pokersLeft = []poker.IPoker{}
	p.pokersPlayed = []poker.IPoker{}
	p.status = Sit
}

// 开始监听用来发来的消息
func (p *Player) startListening() {
	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("player read message error: ", err)
			log.Println("player break")
			p.status = Break
			break
		}

		log.Printf("player recv: %s", msg)

		_msg, _msgErr := message.Decode(msg)
		_msg.PlayerCurrent = p.I

		if _msgErr != nil {
			errorMsg := message.GenMessageNoticeError("pleyer 解析消息失败: " + _msgErr.Error())
			p.conn.WriteMessage(websocket.TextMessage, message.Encode(errorMsg))
			continue
		}

		p.send <- _msg
	}
}

// SetRevc 玩家设置接受牌桌信息的管道
// 设置完立刻监听
func (p *Player) SetRevc(recv chan message.Message) {
	p.recv = recv
	go func() {
		for {
			select {
			case msg := <-p.recv:
				if p.status == Break {
					fmt.Println(strconv.Itoa(p.I) + "号玩家 recv: " + msg.String())
					fmt.Println(strconv.Itoa(p.I) + "号玩家已掉线, 消息丢失")
					continue
				}
				p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
				fmt.Println(strconv.Itoa(p.I) + "号玩家 recv: " + msg.String())
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
