package table

import (
	"ddz/message"
	"ddz/player"
	"ddz/poker"
	"ddz/table/ruler"
	"math/rand"
	"time"

	"errors"
)

type Status = int

const (
	Prepare Status = 0 // 准备中
	Playing Status = 1 // 玩游戏中
	End     Status = 2 // 结束
)

type StatusPlaying = int

const (
	PlayingAllReady = 0 // 全部已准备
	PlayingShuffled = 1 // 已洗牌
	PlayingRealed   = 2 // 已发牌
	PlayingPlaying  = 3 // 出牌中
)

type ITable interface {
	Players() []player.IPlayer           // 列出所有玩家
	PlayerSit(int, player.IPlayer) error // 玩家坐下
}

// Table 桌子
type Table struct {
	i            int                    // 桌号
	status       Status                 // 状态
	subStatus    StatusPlaying          // playing子状态
	players      []player.IPlayer       // 玩家
	pokers       []poker.IPoker         // 牌桌上的牌
	ruler        ruler.IRuler           // 规则检查器
	recvChannel  chan message.Message   // 接受频道
	sendChannels []chan message.Message // 发送频道
	full         int

	playerCurrent int // 当前出牌人
}

// NewTable 新建一个牌桌
func NewTable(i int) *Table {
	t := &Table{
		i:           i,
		status:      Prepare,
		players:     []player.IPlayer{nil, nil, nil, nil},
		pokers:      poker.OnePack(),
		ruler:       ruler.NewRuler(),
		full:        4,
		recvChannel: make(chan message.Message, 1000),
		sendChannels: []chan message.Message{
			make(chan message.Message, 10),
			make(chan message.Message, 10),
			make(chan message.Message, 10),
			make(chan message.Message, 10),
		},
	}
	return t
}

// Players 列出所有玩家
func (p *Table) Players() []player.IPlayer {
	return p.players
}

// PlayerSit 玩家指定一个位置坐下
func (p *Table) PlayerSit(position int, player player.IPlayer) error {
	if p.players[position] != nil {
		return errors.New("该位置已有人")
	}

	p.players[position] = player
	player.Sit(position)
	player.SetRevc(p.sendChannels[position])
	player.SetSend(p.recvChannel)

	return nil
}

// 洗牌
func (p *Table) shuffle() {
	rand.Seed(time.Now().Unix())
	max := len(p.players) - 1
	for i := 0; i <= 1000; i++ {
		fst := rand.Intn(max)
		snd := rand.Intn(max)
		p.pokers[fst], p.pokers[snd] = p.pokers[snd], p.pokers[fst]
	}
}

// 发牌
func (p *Table) real() {
	for i := 0; i < p.full; i++ {
		min := i * 13
		max := (i + 1) * 13
		p.sendChannels[i] <- message.Message{
			T:             message.TypeRuler,
			ST:            message.SubTypeRulerReal,
			Chat:          "",
			PlayerCurrent: i,
			Pokers:        p.pokers[min:max],
		}
		p.pokers = p.pokers[max:]
	}
}

// DaemonRun 后台定时执行
func (p *Table) DaemonRun() {
	go func() {
		// 每秒检查一次事件变动
		t := time.NewTicker(time.Second)
		for {
			<-t.C
			p.daemonCheck()
		}
	}()

	go func() {
		// 检查消息
		p.daemonRecv()
	}()
}

// 检查事件变动
func (p *Table) daemonCheck() {
	//for k, p := range p.Players() {

	//}
}

// 检查消息
func (p *Table) daemonRecv() {
	for {
		msg := <-p.recvChannel
		if msg.T == message.TypeChat {
			for k := range p.sendChannels {
				go func(k int) {
					p.sendChannels[k] <- msg
				}(k)
			}
		}
	}
}
