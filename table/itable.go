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

	// 设置玩家通信管道
	player.SetRevc(p.sendChannels[position])
	player.SetSend(p.recvChannel)

	// 设置玩家在牌桌上的位置
	p.players[position] = player

	// 通知玩家就坐玩家就坐
	player.Sit(position)

	return nil
}

// 洗牌
func (p *Table) shuffle() {
	max := len(p.pokers) - 1
	for i := 0; i <= 1000; i++ {
		rand.Seed(time.Now().UnixNano())
		fst := rand.Intn(max + 1)
		rand.Seed(time.Now().UnixNano() + int64(fst))
		snd := rand.Intn(max + 1)
		p.pokers[fst], p.pokers[snd] = p.pokers[snd], p.pokers[fst]
	}

	p.broadcast(message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerShuffle,
		Chat:          "",
		PlayerCurrent: 0,
		Pokers:        nil,
	})
}

// 发牌
func (p *Table) real() {
	for i := 0; i < p.full; i++ {
		p.sendChannels[i] <- message.Message{
			T:             message.TypeRuler,
			ST:            message.SubTypeRulerReal,
			Chat:          "",
			PlayerCurrent: i,
			Pokers:        p.pokers[0:13],
		}
		p.pokers = p.pokers[13:]
	}
}

// DaemonRun 后台定时执行
func (p *Table) DaemonRun() {
	go func() {
		for {
			msg := <-p.recvChannel
			switch msg.T {
			// 聊天信息
			case message.TypeChat:
				p.broadcast(msg)
			// 游戏中
			case message.TypeRuler:
				switch msg.ST {
				// 已就坐
				case message.SubTypeRulerSit:
					p.broadcast(msg)
				// 已准备
				case message.SubTypeRulerReady:
					p.broadcast(msg)
					// 每次收到玩家已准备信息都检查一次是否可以开始打牌了
					if p.allReady() {
						// 洗牌
						p.shuffle()

						// 停顿1秒，给客户端显示洗牌画面
						time.Sleep(time.Second)

						// 发牌
						p.real()
					}
				}
			}
		}
	}()
}

// 广播信息
func (p *Table) broadcast(msg message.Message) {
	for k := range p.sendChannels {
		go func(k int) {
			p.sendChannels[k] <- msg
		}(k)
	}
}

func (p *Table) sendone(i int, msg message.Message) {
	go func() {
		p.sendChannels[i] <- msg
	}()
}

// 检查是否所有玩家都准备好了
func (p *Table) allReady() bool {
	for _, v := range p.players {
		if v == nil {
			return false
		}
		if v.Status() != player.Prepare {
			return false
		}
	}
	return true
}
