package table

import (
	"ddz/game/message"
	"ddz/game/player"
	"ddz/game/poker"
	"ddz/game/proc1"
	"ddz/game/table/ruler"

	"errors"
	"time"
)

// Status 状态
type Status = int

const (
	// Prepare 准备中
	Prepare Status = 0
	// Playing 玩游戏中
	Playing Status = 1
	// End 结束
	End Status = 2
)

// StatusPlaying 玩牌中的子状态
type StatusPlaying = int

const (
	// PlayingAllReady 全部已准备
	PlayingAllReady StatusPlaying = 0
	// PlayingShuffled 已洗牌
	PlayingShuffled StatusPlaying = 1
	// PlayingRealed 已发牌
	PlayingRealed StatusPlaying = 2
	// PlayingPlaying 出牌中
	PlayingPlaying StatusPlaying = 3
)

// ITable 牌桌
type ITable interface {
	Players() []player.IPlayer           // 列出所有玩家
	PlayerSit(int, player.IPlayer) error // 玩家坐下
}

// Table 桌子
type Table struct {
	i               int                              // 桌号
	status          Status                           // 状态
	subStatus       StatusPlaying                    // playing子状态
	players         []player.IPlayer                 // 玩家
	pokers          []poker.IPoker                   // 牌桌上的牌
	ruler           ruler.IRuler                     // 规则检查器
	recvChannel     chan message.Message             // 接受频道
	sendChannels    []chan message.Message           // 发送频道
	full            int                              // 牌桌满用户数
	maxPokers       []poker.IPoker                   // 当前最大牌
	playerMaxPokers int                              // 当前最大牌的出牌玩家
	playerCurrent   int                              // 当前出牌人
	processors      map[message.Type]proc1.Processor // 消息处理器
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
			make(chan message.Message, 100),
			make(chan message.Message, 100),
			make(chan message.Message, 100),
			make(chan message.Message, 100),
		},
		processors: make(map[message.Type]proc1.Processor),
	}

	setter, _ := t.getSetProcessor()
	setter(message.TypeChat, proc4Chat(t))
	setter(message.TypeRuler, mw4PlayIsMyTurn(t)(proc4Ruler(t)))

	return t
}

// DaemonRun 后台定时执行
func (p *Table) DaemonRun() {
	_, getter := p.getSetProcessor()
	go func() {
		for {
			msg := <-p.recvChannel
			getter(msg.T).Process(msg)
		}
	}()
}

// Players 列出所有玩家
func (p *Table) Players() []player.IPlayer {
	return p.players
}

// 设置或者获取处理器
func (p *Table) getSetProcessor() (
	func(message.Type, proc1.Processor), func(message.Type) proc1.Processor) {

	var setter = func(t message.Type, f proc1.Processor) {
		p.processors[t] = f
	}

	var getter = func(t message.Type) proc1.Processor {
		return p.processors[t]
	}
	return setter, getter
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
	p.ruler.Shuffle(p.pokers)
	p.broadcast(message.GenMessageShuffle())
}

// 发牌
func (p *Table) real() {
	for i := 0; i < p.full; i++ {
		p.sendone(i, message.GenMessageReal(i, p.pokers[0:13]))
		p.Players()[i].SetLeft(p.pokers[0:13])
		p.pokers = p.pokers[13:]
	}
}

// 全都准备好了吗?
func (p *Table) allReady() bool {
	for _, v := range p.players {
		if v == nil || v.Status() != player.Ready {
			return false
		}
	}
	return true
}

// 准备好开始打牌啦
func (p *Table) ready() {
	if p.allReady() {
		p.shuffle()
		time.Sleep(time.Second)
		p.real()
		p.nextPlayer(-1)
	}
}

//切换到下一个用户
func (p *Table) nextPlayer(current int) {
	var _nextPlayer = func(current int) int {
		return (current + 1) % len(p.players)
	}(current)
	p.playerCurrent = _nextPlayer
	p.broadcast(message.GenMessageChangePlayer(_nextPlayer))
}

func (p *Table) broadcast(msg message.Message) {
	for k := range p.sendChannels {
		p.sendone(k, msg)
	}
}

func (p *Table) sendone(playerIndex int, msg message.Message) {
	msg.Send(p.sendChannels[playerIndex])
}

func (p *Table) end(winner int) {
	var restart = func() {
		p.status = Prepare
		p.playerCurrent = -1
		p.pokers = poker.OnePack()

		for i := 0; i < p.full; i++ {
			p.Players()[i].Restart()
		}
	}
	p.broadcast(message.GenMessageWinner(winner))
	p.broadcast(message.GenMessageEnd())
	restart()
}
