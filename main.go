package main

import (
	"ddz/player"
	"ddz/room"
	"ddz/table"
	"fmt"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	// 模拟一些玩家
	player0 := player.NewPlayer("零桌玩家1")
	player1 := player.NewPlayer("零桌玩家2")
	player2 := player.NewPlayer("零桌玩家3")
	player3 := player.NewPlayer("零桌玩家4")
	player4 := player.NewPlayer("一桌玩家1")
	player5 := player.NewPlayer("一桌玩家2")
	player6 := player.NewPlayer("一桌玩家3")

	// 模拟一些桌子
	table0 := table.NewTable(0)
	table1 := table.NewTable(0)
	table2 := table.NewTable(0)
	table3 := table.NewTable(0)
	table4 := table.NewTable(0)
	table5 := table.NewTable(0)

	// 模拟玩家坐到桌位
	table0.PlayerSit(0, player0)
	table0.PlayerSit(1, player1)
	table0.PlayerSit(2, player2)
	table0.PlayerSit(3, player3)
	table1.PlayerSit(0, player4)
	table1.PlayerSit(1, player5)
	table1.PlayerSit(1, player6)

	// 模拟房间，放入6个桌子
	room1 := room.NewRoom([]table.ITable{
		table0, table1, table2,
		table3, table4, table5,
	})

	//fmt.Println(room1.Tables())

	for kt, t := range room1.Tables() {
		tname := "桌" + strconv.Itoa(kt) + ":"
		fmt.Println(tname)
		for kp, p := range t.Players() {
			pname := ""
			if p == nil {
				pname = "没有玩家"
			} else {
				pname = p.Name()
			}
			fmt.Println("第" + strconv.Itoa(kp) + "个位子: " + pname)
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
