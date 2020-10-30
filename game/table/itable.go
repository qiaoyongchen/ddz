package table

import (
	"ddz/game/player"
	"ddz/game/poker"
	"ddz/game/table/ruler"
	"ddz/message"

	"errors"
	"math/rand"
	"time"
)

type Status = int

const (
	Prepare Status = 0 // 准备中
	Playing Status = 1 // 玩游戏中
	End     Status = 2 // 结束
)

// StatusPlaying 玩牌中的子状态
type StatusPlaying = int

const (
	PlayingAllReady StatusPlaying = 0 // 全部已准备
	PlayingShuffled StatusPlaying = 1 // 已洗牌
	PlayingRealed   StatusPlaying = 2 // 已发牌
	PlayingPlaying  StatusPlaying = 3 // 出牌中
)

// ITable 牌桌
type ITable interface {
	Players() []player.IPlayer           // 列出所有玩家
	PlayerSit(int, player.IPlayer) error // 玩家坐下
}

// Table 桌子
type Table struct {
	i               int                        // 桌号
	status          Status                     // 状态
	subStatus       StatusPlaying              // playing子状态
	players         []player.IPlayer           // 玩家
	pokers          []poker.IPoker             // 牌桌上的牌
	ruler           ruler.IRuler               // 规则检查器
	recvChannel     chan message.Message       // 接受频道
	sendChannels    []chan message.Message     // 发送频道
	full            int                        // 牌桌满用户数
	maxPokers       []poker.IPoker             // 当前最大牌
	playerMaxPokers int                        // 当前最大牌的出牌玩家
	playerCurrent   int                        // 当前出牌人
	processors      map[message.Type]processor // 消息处理器
}

// NewTable 新建一个牌桌
func NewTable(i int) *Table {
	t := &Table{
		i:             i,
		status:        Prepare,
		players:       []player.IPlayer{nil, nil, nil, nil},
		pokers:        poker.OnePack(),
		ruler:         ruler.NewRuler(),
		full:          4,
		recvChannel:   make(chan message.Message, 1000),
		playerCurrent: -1,
		sendChannels: []chan message.Message{
			make(chan message.Message, 10),
			make(chan message.Message, 10),
			make(chan message.Message, 10),
			make(chan message.Message, 10),
		},
		processors: make(map[message.Type]processor),
	}

	setter, _ := t.getAndSetProcessor()
	setter(message.TypeChat, chatProcessor)
	setter(message.TypeRuler, forPlayIsMyTurn(rulerProcessor))
	return t
}

// Players 列出所有玩家
func (p *Table) Players() []player.IPlayer {
	return p.players
}

func (p *Table) getAndSetProcessor() (
	setter func(message.Type, processorFunc), getter func(message.Type) processor) {
	return func(t message.Type, f processorFunc) {
			p.processors[t] = f
		}, func(t message.Type) processor {
			return p.processors[t]
		}
}

// PlayerSit 玩家指定一个位置坐下
func (p *Table) PlayerSit(position int, player player.IPlayer) error {
	if p.players[position] != nil {
		return errors.New("该位置已有人")
	}
	player.SetRevc(p.sendChannels[position])
	player.SetSend(p.recvChannel)
	p.players[position] = player
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

// 全都准备好了吗?
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

// 准备好开始打牌啦
// 洗牌 -> 发牌 ->  指定第一个出牌玩家
func (p *Table) ready() {
	if p.allReady() {
		p.shuffle()
		time.Sleep(time.Second)
		p.real()
		p.nextPlayer(-1)
	}
}

// DaemonRun 后台定时执行
func (p *Table) DaemonRun() {
	_, getter := p.getAndSetProcessor()
	go func() {
		for {
			msg := <-p.recvChannel
			getter(msg.T).process(p, msg)
		}
	}()
}

// NextPlayer 切换到下一个用户
func (p *Table) nextPlayer(current int) {
	nextplayer := p._nextPlayer(current)
	p.broadcast(message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerChangePlayer,
		Chat:          "",
		PlayerCurrent: nextplayer,
		Pokers:        nil,
	})
	p.playerCurrent = nextplayer
}

// 获取下一个用户
func (p *Table) _nextPlayer(current int) int {
	return (current + 1) % len(p.players)
}

// 广播信息
func (p *Table) broadcast(msg message.Message) {
	for k := range p.sendChannels {
		p.sendChannels[k] <- msg
	}
}

// 发送给单个人
func (p *Table) sendone(i int, msg message.Message) {
	p.sendChannels[i] <- msg
}

// 广播游戏结束通知
func (p *Table) end(winner int) {
	p.broadcast(message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerWinner,
		Chat:          "",
		PlayerCurrent: winner,
		PlayerTurn:    0,
		Pokers:        nil,
	})
	p.broadcast(message.Message{
		T:             message.TypeRuler,
		ST:            message.SubTypeRulerEnd,
		Chat:          "",
		PlayerCurrent: -1,
		Pokers:        nil,
	})
	// 重新待玩家准备
	p.status = Prepare
	p.playerCurrent = -1
	p.pokers = poker.OnePack()
}
