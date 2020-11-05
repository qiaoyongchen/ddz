package web

import (
	"context"
	"ddz/game"
	"ddz/game/message"
	"ddz/game/player"
	"fmt"
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
		fmt.Println("upgrade error: ", err)
		ctx.JSON(http.StatusOK, "upgrade error")
		return nil
	}

	// -----------------------------自动检测断线重连开始-------------------------------------
	// TODO 这块业务逻辑需要重写先暂时这么实现
	_room := game.GetRoom()
	tables := _room.Tables()
	for _, _table := range tables {
		players := _table.Players()
		for _, _player := range players {
			if _player != nil && _player.IfBreak() {
				relinkErr := _player.RelinkWhenBreaking(conn)
				if relinkErr != nil {
					conn.WriteMessage(websocket.TextMessage,
						message.Encode(message.GenMessageNoticeError(relinkErr.Error())))
				} else {
					return nil
				}
			}
		}
	}
	// -----------------------------自动检测断线重连结束-------------------------------------

	conn.WriteMessage(websocket.TextMessage,
		message.Encode(message.GenMessageRoomInfo(game.GetRoomInfo())))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message: ", err)
			ctx.JSON(http.StatusOK, "read message error")
			break
		}
		fmt.Println("recv: ", msg)

		_msg, _msgErr := message.Decode(msg)
		if _msgErr != nil {
			conn.WriteMessage(websocket.TextMessage,
				message.Encode(message.GenMessageNoticeError("解析消息失败: "+_msgErr.Error())))
			continue
		}

		sitted := false

		switch _msg.T {
		case message.TypeRuler:
			switch _msg.ST {
			case message.SubTypeRulerSit:
				p := player.NewPlayer("123", conn)
				_room := game.GetRoom()
				_room.Tables()[_msg.TableIndex].PlayerSit(_msg.TablePositionIndex, p)
				sitted = true
				goto ENDING
			default:
				goto ENDING
			}
		case message.TypeRoom:
			switch _msg.ST {
			case message.SubTypeGetRoomInfo:
				conn.WriteMessage(websocket.TextMessage,
					message.Encode(message.GenMessageNoticeError("解析消息失败: "+_msgErr.Error())))
				goto ENDING
			default:
				goto ENDING
			}
		default:
			goto ENDING
		}

	ENDING:
		if sitted == true {
			fmt.Println("conn 交给用户管理 ")
			break
		}
	}
	return nil
}
