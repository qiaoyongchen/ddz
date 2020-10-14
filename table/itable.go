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

type Status int

const (
	Prepare Status = 0 // 准备中
	Playing Status = 1 // 玩游戏中
	End     Status = 2 // 结束
)

type ITable interface {
	Players() []player.IPlayer           // 列出所有玩家
	PlayerSit(int, player.IPlayer) error // 玩家坐下
}

// Table 桌子
type Table struct {
	i            int                    // 桌号
	status       Status                 // 状态
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
		recvChannel: make(chan message.Message),
		sendChannels: []chan message.Message{
			make(chan message.Message),
			make(chan message.Message),
			make(chan message.Message),
			make(chan message.Message),
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

}
