package player

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"

	"ddz/game/poker"
	"ddz/message"
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
	Prepare() error               // 准备
	Play([]poker.IPoker) error    // 出牌
	SetRevc(chan message.Message) // 设置接受管道
	SetSend(chan message.Message) // 设置发送管道
	GetName() string              // 玩家名字
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
	conn         *websocket.Conn      // websocket 连接
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
	p.I = i
	p.status = Sit
	p.startListening()
	p.send <- message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerSit,
		PlayerCurrent: i,
		PlayerTurn:    0,
	}
}

// Prepare 准备
func (p *Player) Prepare() error {
	if p.status != Sit {
		return errors.New("状态出错，不能准备开始")
	}
	p.status = Prepare
	return nil
}

// Play 出牌
func (p *Player) Play([]poker.IPoker) error {
	return nil
}

// 开始监听用来发来的消息
func (p *Player) startListening() {
	go func() {
		for {
			_, msg, err := p.conn.ReadMessage()
			if err != nil {
				log.Println("player read message: ", err)
				p.conn.WriteMessage(websocket.TextMessage, []byte("player read message: "+err.Error()))
				continue
			}
			log.Printf("player recv: %s", msg)

			_msg, _msgErr := message.Decode(msg)
			if _msgErr != nil {
				p.conn.WriteMessage(websocket.TextMessage, message.Encode(
					message.Message{
						T:    message.TypeNotice,
						ST:   message.SubTypeNoticeError,
						Chat: "pleyer 解析消息失败: " + _msgErr.Error(),
					},
				))
				continue
			}
			switch _msg.T {
			case message.TypeRuler:
				switch _msg.ST {
				case message.SubTypeRulerReady:
					p.status = Prepare
				default:
					continue
				}
			default:
				continue
			}
			p.send <- _msg
		}
	}()
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
					fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [" + msg.Chat + "]")
					p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
				case message.TypeNotice:
					fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [" + msg.Chat + "]")
					p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
				case message.TypeRuler:
					switch msg.ST {
					case message.SubTypeRulerSit:
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家已就坐]")
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
					case message.SubTypeRulerReady:
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家已准备]")
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
					case message.SubTypeRulerShuffle:
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [洗牌中]")
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
					case message.SubTypeRulerReal:
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [发牌:( " + poker.ShowPokers(msg.Pokers) + " )]")
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
						p.pokersLeft = msg.Pokers
					case message.SubTypeRulerPlay:
						if len(msg.Pokers) == 0 {
							fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家: 不要]")
							p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
							continue
						}
						newleft, err := poker.SubPokers(p.pokersLeft, msg.Pokers)
						if err != nil {
							p.conn.WriteMessage(websocket.TextMessage, []byte("player play message: "+err.Error()))
							continue
						}
						p.pokersLeft = newleft
						p.pokersPlayed = append(p.pokersPlayed, msg.Pokers...)
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家出牌:( " + poker.ShowPokers(msg.Pokers) + " )]")
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
					case message.SubTypeRulerChangePlayer:
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [现在轮到" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家出牌]")
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
					case message.SubTypeRulerWinner:
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [" + strconv.Itoa(msg.PlayerCurrent) + "号位置玩家获胜]")
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
					case message.SubTypeRulerEnd:
						fmt.Println(strconv.Itoa(p.I) + "号玩家 recive: [本局游戏结束]")
						p.Restart()
						p.conn.WriteMessage(websocket.TextMessage, message.Encode(msg))
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
		Chat:          p.Name + " say: " + content,
		PlayerCurrent: p.I,
	}
}

func (p *Player) Restart() {
	p.pokersLeft = []poker.IPoker{}
	p.pokersPlayed = []poker.IPoker{}
	p.status = Sit
}
