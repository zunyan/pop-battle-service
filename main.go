package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

const (
	ROOM_STATUS_PENDING = 0

	// 用户状态
	PLAYER_STATUS_PENDING = 0
	PLAYER_STATUS_READY   = 1
)

type Player struct {
	Name      string `json:"name"`
	RoleIndex int    `json:"roleIndex"`
	Status    int    `json:"status"`
}
type Room struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Players     []*Player `json:"players"`
	TotalPlayer int       `json:"totalPlayer"`
	Status      int       `json:"status"`
}

type SocketACK bool

func main() {
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
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {})

	//大厅
	roomList := []*Room{}
	roomListMap := map[string]*Room{}

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
	server.OnEvent("/hall", "getRoomList", func(s socketio.Conn, msg string) []*Room {
		return roomList
	})

	// 创建房间
	server.OnEvent("/hall", "createRoom", func(s socketio.Conn, roomName string) Room {
		url := s.URL()
		urlQuery := url.Query()
		username := urlQuery.Get("username")
		room := Room{
			Id:   uuid.New().String(),
			Name: roomName,
			Players: []*Player{
				&Player{Name: username},
			},
			TotalPlayer: 4,
			Status:      ROOM_STATUS_PENDING,
		}
		roomList = append(roomList, &room)
		roomListMap[room.Id] = &room
		fmt.Println(username + "创建了房间" + roomName)
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
		room, exist := roomListMap[roomId]
		fmt.Println(url, username, roomId)
		if !exist {
			return errors.New("无效的房间号")
		}

		s.Join(roomId)
		server.BroadcastToRoom("/room", roomId, "message", username+"进入了房间 "+room.Name)
		server.BroadcastToRoom("/room", roomId, "sync", roomListMap[roomId])
		return nil
	})

	// 选择角色
	server.OnEvent("/room", "choosePlayer", func(s socketio.Conn, roleIndex int) error {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")

		var pleyer *Player
		for _, temp := range roomListMap[roomId].Players {
			if temp.Name == username {
				pleyer = temp
				break
			}
		}

		if pleyer.Name == "" {
			return errors.New("无效的用户")
		}

		pleyer.RoleIndex = roleIndex
		server.BroadcastToRoom("/room", roomId, "sync", roomListMap[roomId])

		return nil
	})

	// 准备
	server.OnEvent("/room", "ready", func(s socketio.Conn, roleIndex int) error {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")

		var pleyer *Player
		for _, temp := range roomListMap[roomId].Players {
			if temp.Name == username {
				pleyer = temp
				break
			}
		}

		if pleyer.Name == "" {
			return errors.New("无效的用户")
		}

		pleyer.Status = PLAYER_STATUS_READY
		server.BroadcastToRoom("/room", roomId, "sync", roomListMap[roomId])

		return nil
	})

	server.OnDisconnect("/room", func(s socketio.Conn, reason string) {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")
		room, exist := roomListMap[roomId]
		if exist {
			roomName := room.Name
			s.Leave(roomId)
			server.BroadcastToRoom("/room", roomId, "sync", roomListMap[roomId])
			server.BroadcastToRoom("/room", roomId, "message", username+"离开了房间"+roomName)
			fmt.Println(username + "离开了房间" + roomName)
			newPlayers := []*Player{}
			for _, v := range room.Players {
				if v.Name != username {
					newPlayers = append(newPlayers, v)
				}
			}

			room.Players = newPlayers

			if len(room.Players) == 0 {
				fmt.Println(roomName, "房间人数不足，自动解散")
				// clear room
			} else {
				server.BroadcastToRoom("/room", roomId, "sync", roomListMap[roomId])
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
