package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"pop-battle-service/pkg/store"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

type SocketACK bool

func main() {

	store.Init()
	logger := log.New(os.Stdout, "<main>", log.Lshortfile|log.Ldate|log.Ltime)
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	// 默认连接, 不可删除
	server.OnConnect("/", func(s socketio.Conn) error {
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		s.Close()
		logger.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {})

	//大厅
	// 进入大厅时进行处理
	server.OnConnect("/hall", func(s socketio.Conn) error {
		url := s.URL()
		urlQuery := url.Query()
		username := urlQuery.Get("username")
		if username == "" {
			return errors.New("请输入用户名")
		}
		logger.Println(username, "加入大厅")
		server.BroadcastToNamespace("/hall", "message", username+"加入大厅")
		return nil
	})

	// 离开大厅
	server.OnDisconnect("/hall", func(s socketio.Conn, reason string) {
		url := s.URL()
		urlQuery := url.Query()
		username := urlQuery.Get("username")
		logger.Println(username, "离开大厅")
		server.BroadcastToNamespace("/hall", "message", username+"离开大厅")
	})

	// 获取房间列表
	server.OnEvent("/hall", "getRoomList", func(s socketio.Conn, msg string) []*store.Room {
		return store.RoomList
	})

	// 创建房间
	server.OnEvent("/hall", "createRoom", func(s socketio.Conn, roomName string) *store.Room {
		url := s.URL()
		urlQuery := url.Query()
		username := urlQuery.Get("username")

		room := store.CreateRoom(roomName, username)
		logger.Println(username + "创建了房间" + roomName)
		server.BroadcastToNamespace("/hall", "message", username+"创建了房间"+roomName)
		return room
	})

	// 房间内操作
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
		logger.Println(username, "加入房间，当前房间人数", len(room.Players))
		server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])
		return nil
	})

	// 选择角色
	server.OnEvent("/room", "choosePlayer", func(s socketio.Conn, role string) error {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")
		var pleyer *store.Player
		for _, temp := range store.RoomListMap[roomId].Players {
			if temp.Name == username {
				pleyer = temp
				break
			}
		}

		if pleyer.Name == "" {
			return errors.New("无效的用户")
		}
		logger.Printf("房间:%v,用户:%v更换角色,由%v更换为%v\n", store.RoomListMap[roomId].Name, username, pleyer.Role, role)
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

		var pleyer *store.Player
		for _, temp := range store.RoomListMap[roomId].Players {
			if temp.Name == username {
				pleyer = temp
				break
			}
		}

		if pleyer.Name == "" {
			return errors.New("无效的用户")
		}

		pleyer.Status = store.PLAYER_STATUS_READY
		server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])

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
			logger.Println(username + "离开了房间" + roomName)
			newPlayers := []*store.Player{}
			for _, v := range room.Players {
				if v.Name != username {
					newPlayers = append(newPlayers, v)
				}
			}

			room.Players = newPlayers

			if len(room.Players) == 0 {
				logger.Println(roomName, "房间人数不足，自动解散")
				// clear room
				newRooms := []*store.Room{}
				for _, v := range store.RoomList {
					if v != room {
						newRooms = append(newRooms, v)
					}
				}
				store.RoomList = newRooms
				delete(store.RoomListMap, roomId)
			} else {
				server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])
			}
		}
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000 ...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
