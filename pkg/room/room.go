package room

import (
	"errors"
	"fmt"
	"pop-battle-service/pkg/store"
	"pop-battle-service/pkg/typings"

	socketio "github.com/googollee/go-socket.io"
)

func LinkRouter(server *socketio.Server) {
	// 加入
	server.OnConnect("/room", func(s socketio.Conn) error {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")

		room, err := store.JoinRoom(roomId, username)
		if err != nil {
			s.Leave("/room")
			s.Close()
			return err
		}

		s.Join(roomId)
		server.BroadcastToRoom("/room", roomId, "message", username+"进入了房间 "+room.Name)
		fmt.Println(username, "加入房间，当前房间人数", len(room.Players))
		server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])
		return nil
	})

	// 选择角色
	server.OnEvent("/room", "choosePlayer", func(s socketio.Conn, role string) error {

		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")
		fmt.Println(username, "更换角色")
		var pleyer *typings.Player
		for _, temp := range store.RoomListMap[roomId].Players {
			if temp.Name == username {
				pleyer = temp
				break
			}
		}

		if pleyer.Name == "" {
			return errors.New("无效的用户")
		}
		fmt.Println(username, pleyer.Role, role)
		pleyer.Role = role
		server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])

		return nil
	})

	// 准备
	server.OnEvent("/room", "ready", func(s socketio.Conn, role string) error {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")

		var pleyer *typings.Player
		for _, temp := range store.RoomListMap[roomId].Players {
			if temp.Name == username {
				pleyer = temp
				break
			}
		}

		if pleyer.Name == "" {
			return errors.New("无效的用户")
		}

		if pleyer.Status == typings.PLAYER_STATUS_READY {
			pleyer.Status = typings.PLAYER_STATUS_PENDING
		} else {
			pleyer.Status = typings.PLAYER_STATUS_READY
		}
		server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])
		return nil
	})

	server.OnEvent("/room", "start", func(s socketio.Conn, role string) error {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		if !store.HasRoom(roomId) {
			return errors.New("房间不存在")
		}

		room := store.RoomListMap[roomId]
		if len(room.Players) == 1 {
			return errors.New("人数不足无法开始")
		}

		fmt.Println("开始游戏", room.Name)
		room.Status = typings.ROOM_STATUS_IN_GAME
		server.BroadcastToRoom("/room", roomId, "gameStart", room)
		return nil
	})

	server.OnDisconnect("/room", func(s socketio.Conn, reason string) {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")
		room, exist := store.RoomListMap[roomId]
		if exist {
			roomName := room.Name
			s.Leave(roomId)
			server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])
			server.BroadcastToRoom("/room", roomId, "message", username+"离开了房间"+roomName)
			fmt.Println(username + "离开了房间" + roomName)
			newPlayers := []*typings.Player{}
			var leavePlayer *typings.Player
			for _, v := range room.Players {
				if v.Name != username {
					newPlayers = append(newPlayers, v)
				} else {
					leavePlayer = v
				}
			}

			room.Players = newPlayers

			if len(room.Players) == 0 {
				fmt.Println(roomName, "房间人数不足，自动解散")
				// clear room
				newRooms := []*typings.Room{}
				for _, v := range store.RoomList {
					if v != room {
						newRooms = append(newRooms, v)
					}
				}
				store.RoomList = newRooms
				delete(store.RoomListMap, roomId)

				return
			}

			if leavePlayer != nil && leavePlayer.IsMaster {
				// 房主离开，自动票选下一个人作为房主
				room.Players[0].IsMaster = true
			}

			server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])
		}
	})
}
