package player

import(
	"ddz/poker"
)

type Status uint8

const (
	NotPrepare = 0 // 未准备
	Prepare = 1 // 已准备 
	Playing= 2 // 进行中
)

type IPlayer interface {
	Left() poker.IPoker // 剩余牌
	Status() Status // 状态 
	Play([]poker.IPoker) error // 出牌
}