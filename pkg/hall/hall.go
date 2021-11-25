package hall

import (
	"errors"
	"fmt"
	"pop-battle-service/pkg/store"
	"pop-battle-service/pkg/typings"

	socketio "github.com/googollee/go-socket.io"
)

func LinkRouter(server *socketio.Server) {
	// 进入大厅时进行处理
	server.OnConnect("/hall", func(s socketio.Conn) error {
		url := s.URL()
		urlQuery := url.Query()
		username := urlQuery.Get("username")
		if username == "" {
			return errors.New("请输入用户名")
		}
		fmt.Println(username, "加入大厅")
		server.BroadcastToNamespace("/hall", "message", username+"加入大厅")
		return nil
	})

	// 离开大厅
	server.OnDisconnect("/hall", func(s socketio.Conn, reason string) {
		url := s.URL()
		urlQuery := url.Query()
		username := urlQuery.Get("username")
		fmt.Println(username, "离开大厅")
		server.BroadcastToNamespace("/hall", "message", username+"离开大厅")
	})

	// 获取房间列表
	server.OnEvent("/hall", "getRoomList", func(s socketio.Conn, msg string) []*typings.Room {
		return store.RoomList
	})

	// 创建房间
	server.OnEvent("/hall", "createRoom", func(s socketio.Conn, roomName string) *typings.Room {
		url := s.URL()
		urlQuery := url.Query()
		username := urlQuery.Get("username")

		room := store.CreateRoom(roomName, username)
		fmt.Println(username + "创建了房间" + roomName)
		server.BroadcastToNamespace("/hall", "message", username+"创建了房间"+roomName)
		return room
	})
}
