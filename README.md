# ddz
```
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
	table1 := table.NewTable(1)
	table2 := table.NewTable(2)

	// 模拟玩家坐到桌位
	table0.PlayerSit(0, player0)
	table0.PlayerSit(1, player1)
	table0.PlayerSit(2, player2)
	table0.PlayerSit(3, player3)
	table1.PlayerSit(0, player4)
	table1.PlayerSit(1, player5)
	table1.PlayerSit(1, player6)

	// 模拟房间，放入3个牌桌
	room1 := room.NewRoom([]table.ITable{table0, table1, table2})

	for kt, t := range room1.Tables() {
		fmt.Println("桌号" + strconv.Itoa(kt) + ":")
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

	// 第一张桌子开张
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

	// 玩家1 出第一张牌
	player0.PlayFirst()
	time.Sleep(time.Second * 1)
	player0.ShowLeft()
	// 玩家1 不出
	player1.PlayNone()
	time.Sleep(time.Second * 1)
	// 玩家2 全出
	player2.PlayAll()
	time.Sleep(time.Second * 1)

	fmt.Println("第一局结束........................................")
	// 结束
	time.Sleep(time.Second * 2)

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

	// 玩家1 出第一张牌
	player0.PlayFirst()
	time.Sleep(time.Second * 2)
	player0.ShowLeft()
	// 玩家1 不出
	player1.PlayNone()
	// 玩家2 全出
	player2.PlayAll()

	// web 服务开启
	time.Sleep(time.Second * 1)
```


-- 123

```
// PlayNone 不出
func (p *Player) PlayNone() error {
	return p.Play([]poker.IPoker{})
}

// PlayAll 全部出掉
func (p *Player) PlayAll() error {
	return p.Play(p.pokersLeft)
}

// PlayFirst 出第一张牌(测试)
func (p *Player) PlayFirst() error {
	return p.Play(p.pokersLeft[0:1])
}
```