package game

import (
	"fmt"
	"strconv"

	"ddz/game/room"
	"ddz/game/table"
)

// 一个大厅放三张牌桌
var (
	table0 = table.NewTable(0)
	table1 = table.NewTable(1)
	table2 = table.NewTable(2)
	room1  = room.NewRoom([]table.ITable{table0, table1, table2})
)

// Run 运行
func Run() {
	table0.DaemonRun()
	table1.DaemonRun()
	table2.DaemonRun()
	println("room ready, tables ready, game start")
}

// Shutdown 退出
func Shutdown() {
	fmt.Println("牌桌0清空...")
	fmt.Println("牌桌1清空...")
	fmt.Println("牌桌2清空...")
	fmt.Println("大厅关闭...")
}

// GetRommInfo 获取大厅信息
func GetRommInfo() interface{} {
	rst := make(map[string]interface{})
	for k, tb := range room1.Tables() {
		rst[strconv.Itoa(k)] = tb.Players()
	}
	return rst
}
