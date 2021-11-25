package game

import (
	"errors"
	"fmt"
	"pop-battle-service/pkg/gamemap"
	"pop-battle-service/pkg/typings"

	socketio "github.com/googollee/go-socket.io"
)

var gameMap map[string]*typings.TGameInfo

func Init() func(server *socketio.Server) {
	gameMap = map[string]*typings.TGameInfo{}
	return link
}

func link(server *socketio.Server) {
	// 加入
	server.OnConnect("/game", func(s socketio.Conn) error {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")

		game, exist := gameMap[roomId]
		if !exist {
			return errors.New("对局信息不存在")
		}

		s.Join(roomId)
		s.Emit("sync", game)
		fmt.Println(username, "加入对局")
		server.BroadcastToRoom("/game", roomId, "sync", game)
		return nil
	})

	server.OnEvent("/game", "changeprop", func(s socketio.Conn, props typings.TGamePlayer) {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")

		game, exist := gameMap[roomId]
		if !exist {
			return
		}

		for _, p := range game.Players {
			if p.Name == username {
				p.Gridx = props.Gridx
				p.Gridy = props.Gridy
				// p.Bubbles = props.Bubbles
				// p.Speed = props.Speed
				// p.Power = props.Power
				p.X = props.X
				p.Y = props.Y
				p.MoveTarget = props.MoveTarget

			}
		}

		server.BroadcastToRoom("/game", roomId, "sync", game)
	})

}

func CreateGame(room *typings.Room) {
	// 生成地图数据

	boxs, roles := gamemap.GetGameMap()
	players := []*typings.TGamePlayer{}
	i := 0

	for _, p := range room.Players {
		players = append(players, &typings.TGamePlayer{
			Gridx:      roles[i][0],
			Gridy:      roles[i][1],
			Name:       p.Name,
			Status:     typings.TGamePlayerStatus_ALIVE,
			Speed:      4,
			Power:      3,
			Bubbles:    2,
			X:          roles[i][0]*40 + 20,
			Y:          roles[i][1]*40 + 20,
			MoveTarget: typings.TGamePlayerMoveTarget_None,
		})

		i++
	}

	gameInfo := &typings.TGameInfo{
		Props:   boxs,
		Players: players,
	}

	gameMap[room.Id] = gameInfo
}
