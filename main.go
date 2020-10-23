package main

import (
	"ddz/game"
	"ddz/web"
	"os"
	"os/signal"
)

func main() {
	run()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	shutdown()
}

func run() {
	game.Run()
	web.Run()
}

// 退出收尾程序
func shutdown() {
	web.Shutdown()
	game.Shutdown()
}
