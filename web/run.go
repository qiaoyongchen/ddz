package web

import (
	"context"
	"ddz/game"
	"ddz/game/player"
	"ddz/message"
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
		srv.Start(":1234")
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
	conn.WriteMessage(websocket.TextMessage, message.Encode(
		message.Message{
			T:    message.TypeRoom,
			ST:   message.SubTypeRoomInfo,
			Data: game.GetRoomInfo(),
		},
	))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message: ", err)
			ctx.JSON(http.StatusOK, "read message error")
			break
		}
		log.Printf("recv: %s", msg)

		_msg, _msgErr := message.Decode(msg)
		if _msgErr != nil {
			conn.WriteMessage(websocket.TextMessage, message.Encode(
				message.Message{
					T:    message.TypeNotice,
					ST:   message.SubTypeNoticeError,
					Chat: "解析消息失败: " + _msgErr.Error(),
				},
			))
			continue
		}
		switch _msg.T {
		case message.TypeRuler:
			switch _msg.ST {
			case message.SubTypeRulerSit:
				p := player.NewPlayer("123", conn)
				_room := game.GetRoom()
				_room.Tables()[_msg.TableIndex].PlayerSit(_msg.TablePositionIndex, p)
				break
			default:
				continue
			}
		default:
			continue
		}
	}
	return nil
}
