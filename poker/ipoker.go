package poker

type Type = uint

const (
	TypeS = 10000
	TypeB = 10001
	TypeHeiTao = 3 // 黑桃
	TypeHongTao = 2 //红桃
	TypeMeihua = 1 // 梅花
	TypeFangkuai = 0 // 方块
)

type Value = uint

const (
	V3 = 3
	V4 = 4
	V5 = 3
	V6 = 3
	V7 = 3
	V8 = 3
	V9 = 3
	V10 = 3
	VJ = 11
	VQ = 12
	VK = 13
	VA = 14
	V2 = 15
	VS = 10000
	VB = 10001
)

type IPoker interface {
	Type() Type
	Value() Value
}