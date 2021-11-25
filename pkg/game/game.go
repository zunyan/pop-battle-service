package game

import (
	socketio "github.com/googollee/go-socket.io"
)

// TGameProp
type TGameProp uint

const (
	GAME_PROPS_SHOSE   = 1
	GAME_PROPS_LOTION  = 2
	GAME_PROPS_BUBBLES = 3
)

// TGamePlayerStatus
type TGamePlayerStatus = uint

const (
	GAME_PROPS_ALIVE = 1
	GAME_PROPS_DEAD  = 2
)

type TGamePlayer struct {
	Gridx   int               `json:"gridx"`
	Gridy   int               `json:"gridy"`
	Name    string            `json:"name"`
	Status  TGamePlayerStatus `json:"status"`
	Speed   int               `json:"speed"`
	Power   int               `json:"power"`
	Bubbles int               `json:"bubbles"`
}
type TGameInfo struct {
	props   []TGameProp    `json:"props"`
	players []*TGamePlayer `json:"players"`
}

var gameMap map[string]*TGameInfo

func Init() func(server *socketio.Server) {
	gameMap = map[string]*TGameInfo{}
	return link
}

func link(server *socketio.Server) {
	// 加入
	server.OnConnect("/room", func(s socketio.Conn) error {
		// url := s.URL()
		// urlQuery := url.Query()
		// roomId := urlQuery.Get("roomId")
		// username := urlQuery.Get("username")

		// game, exist := gameMap[]
		// room, err := store.JoinRoom(roomId, username)
		// if err != nil {
		// 	s.Leave("/room")
		// 	s.Close()
		// 	return err
		// }

		// s.Join(roomId)
		// server.BroadcastToRoom("/room", roomId, "message", username+"进入了房间 "+room.Name)
		// fmt.Println(username, "加入房间，当前房间人数", len(room.Players))
		// server.BroadcastToRoom("/room", roomId, "sync", store.RoomListMap[roomId])
		return nil
	})
}
