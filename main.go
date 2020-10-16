package main

import (
	"ddz/player"
	"ddz/poker"
	"ddz/room"
	"ddz/table"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	// 模拟一些玩家
	player0 := player.NewPlayer("零桌玩家0")
	player1 := player.NewPlayer("零桌玩家1")
	player2 := player.NewPlayer("零桌玩家2")
	player3 := player.NewPlayer("零桌玩家3")
	player4 := player.NewPlayer("一桌玩家0")
	player5 := player.NewPlayer("一桌玩家1")
	player6 := player.NewPlayer("一桌玩家2")

	// 模拟一些桌子
	table0 := table.NewTable(0)
	table1 := table.NewTable(0)
	table2 := table.NewTable(0)

	// 模拟玩家坐到桌位
	table0.PlayerSit(0, player0)
	table0.PlayerSit(1, player1)
	table0.PlayerSit(2, player2)
	table0.PlayerSit(3, player3)
	table1.PlayerSit(0, player4)
	table1.PlayerSit(1, player5)
	table1.PlayerSit(1, player6)

	// 模拟房间，放入3个桌子
	room1 := room.NewRoom([]table.ITable{table0, table1, table2})

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

	// 第一张桌子开始开张
	table0.DaemonRun()

	// 玩家1进行聊天，测试广播
	println()
	player0.Chat("hello, every one")

	// 玩家准备开始玩游戏
	player0.Prepare()
	time.Sleep(time.Second / 4)
	player1.Prepare()
	time.Sleep(time.Second / 4)
	player2.Prepare()
	time.Sleep(time.Second / 4)
	player3.Prepare()

	// 停顿两秒再出牌,模拟洗牌和发牌过程中的耗时
	time.Sleep(time.Second * 2)

	// 用户0出♥3，该玩家如果有的话正常显示，该玩家如果没有的话则报错
	// 需要多试几次
	player0.Play([]poker.IPoker{poker.NewPoker(poker.TypeHeart, poker.V3)})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
