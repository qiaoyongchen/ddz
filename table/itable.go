package table

type Status int

const (
	Prepare Status = 0 // 准备中
	Playing Status = 1 // 进行中
	End Status = 2 // 结束
)

type ITable interface {

}