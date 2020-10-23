package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var srv = echo.New()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	srv.HideBanner = true
	srv.Use(middleware.CORS())
	srv.Use(middleware.Logger())
	srv.Use(middleware.Recover())
	srv.GET("/helloworld", test)
	srv.GET("/connect", websocketRun)
}

// Run 运行
func Run() {
	go func() {
		srv.Logger.Fatal(srv.Start(":1234"))
	}()
}

// Shutdown 关闭
func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger.Fatal(err)
	}
}

func test(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, World!")
}

func websocketRun(ctx echo.Context) error {
	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Print("upgrade error: ", err)
		ctx.JSON(http.StatusOK, "upgrade error")
		return nil
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message: ", err)
			ctx.JSON(http.StatusOK, "read message error")
			break
		}
		log.Println(msgType)
		log.Printf("recv: %s", msg)
	}
	return nil
}
