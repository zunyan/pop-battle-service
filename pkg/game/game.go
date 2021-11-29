package game

import (
	"errors"
	"fmt"
	"pop-battle-service/pkg/gamemap"
	"pop-battle-service/pkg/typings"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

var gameMap map[string]*Game
var mserver *socketio.Server

func Init() func(server *socketio.Server) {
	gameMap = map[string]*Game{}
	return link
}

func link(server *socketio.Server) {
	mserver = server
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
		fmt.Println(username, "加入对局")
		server.BroadcastToRoom("/game", roomId, "sync", game.getSyncData())
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

		player := game.getPlayerByName(username)
		if player == nil {
			return
		}

		player.Gridx = props.Gridx
		player.Gridy = props.Gridy
		player.X = props.X
		player.Y = props.Y
		player.MoveTarget = props.MoveTarget

		server.BroadcastToRoom("/game", roomId, "sync", game.getSyncData())
	})

	server.OnEvent("/game", "putBubbles", func(s socketio.Conn) {
		url := s.URL()
		urlQuery := url.Query()
		roomId := urlQuery.Get("roomId")
		username := urlQuery.Get("username")

		game, exist := gameMap[roomId]
		if !exist {
			return
		}
		player := game.getPlayerByName(username)
		if player != nil && game.addBubble(player.Gridx, player.Gridy, player.Power) {
			server.BroadcastToRoom("/game", roomId, "sync", game.getSyncData())
		}
	})

}

type TGameBubbles []*typings.TGameBubble

type Game struct {
	RoomId           string
	Props            []*typings.TGameBox    `json:"props"`
	Players          []*typings.TGamePlayer `json:"players"`
	Bubbles          TGameBubbles           `json:"bubbles"`
	GameMap          [][]*gamemap.TGameMapBlock
	CheckThreadExist bool
}

type TGameSyncPack struct {
	Props   []*typings.TGameBox    `json:"props"`
	Players []*typings.TGamePlayer `json:"players"`
	Bubbles TGameBubbles           `json:"bubbles"`
}

func (this *Game) boom(bnb *typings.TGameBubble) []*typings.TGameBoomBubble {
	todoList := []*typings.TGameBubble{bnb}
	booms := []*typings.TGameBoomBubble{}
	destoryBox := []*typings.TGameBox{}

	fn := func(temp *typings.TGameBubble, step int, prop string) int {
		nGridx := temp.Gridx
		nGridy := temp.Gridy
		l := 0

		for {

			if prop == "x" {
				nGridx += step
			} else {
				nGridy += step
			}

			if nGridx >= len(this.GameMap[0]) || nGridx < 0 || nGridy >= len(this.GameMap) || nGridy < 0 {
				break
			}

			nextItem := this.GameMap[nGridy][nGridx]

			if !nextItem.CanDestory {
				break
			}

			// 如果目标位置有炸弹
			if nextItem.Bubble != nil {
				// 放到todolist里面，等待下次检查
				todoList = append(todoList, nextItem.Bubble)

				// 把地图的栅格炸弹标记为nil， 防止下一个循环再次进来
				nextItem.Bubble = nil
			}

			// 如果这个地方有箱子
			// 如果限制没被销毁，进行标记，同时此处应该为水流的最大位置
			// 再循环遍历的过程中，可能多个球会同时命中一个箱子，如果此时就将箱子改为nil， 那么有下一个球就会默认此处没有障碍物，而继续往前判断
			if nextItem.Box != nil {
				if !nextItem.Box.Hasdestoryed {
					nextItem.Box.Hasdestoryed = true
					l++
					destoryBox = append(destoryBox, nextItem.Box)
				}
				break
			}

			l++
			if l >= bnb.Power {
				break
			}
		}
		return l
	}

	i := 0
	this.GameMap[bnb.Gridy][bnb.Gridx].Bubble = nil
	for len(todoList) > i {

		temp := todoList[i]
		right := fn(temp, 1, "x")
		left := fn(temp, -1, "x")
		top := fn(temp, -1, "y")
		bottom := fn(temp, 1, "y")

		temp.Hasdestoryed = true

		booms = append(booms, &typings.TGameBoomBubble{
			Gridx:  temp.Gridx,
			Gridy:  temp.Gridy,
			Left:   left,
			Right:  right,
			Top:    top,
			Bottom: bottom,
		})
		i++
	}

	// 将已经销毁的泡泡从 bubbles 中删除
	newBubbles := []*typings.TGameBubble{}
	for _, v := range this.Bubbles {
		if !v.Hasdestoryed {
			newBubbles = append(newBubbles, v)
		}
	}
	this.Bubbles = newBubbles

	for _, v := range destoryBox {
		this.GameMap[v.Gridy][v.Gridx].Box = nil
	}

	return booms
}

func (game *Game) addBubble(gridX int, gridY int, power int) bool {
	grid := game.GameMap[gridY][gridX]
	if grid.Bubble != nil || grid.Box != nil {
		fmt.Println("已经存在炸弹，或者改位置仍然有障碍物", grid.Bubble)
		return false
	}

	bubble := &typings.TGameBubble{
		Gridx:      gridX,
		Gridy:      gridY,
		Power:      power,
		CreateTime: time.Duration(time.Nanosecond),
	}

	grid.Bubble = bubble
	game.Bubbles = append(game.Bubbles, bubble)

	if !game.CheckThreadExist {
		game.CheckThreadExist = true
		go game.checkBubble()
	}

	return true
}

func (game *Game) checkBubble() {
	roomId := game.RoomId
	for {
		if len(game.Bubbles) == 0 {
			break
		}
		bubble := game.Bubbles[0]
		t := bubble.CreateTime + 2000
		fmt.Println("距离炸弹爆炸事件ms：", t-time.Duration(time.Now().UnixNano()/1e6))
		time.Sleep(2 * time.Second) // 此处有问题

		booms := game.boom(bubble)

		// 发给客户端
		mserver.BroadcastToRoom("/game", roomId, "boomBubble", booms)
		// sync 消息
		mserver.BroadcastToRoom("/game", roomId, "sync", game.getSyncData())

	}

	game.CheckThreadExist = false

}

func (game *Game) getSyncData() *TGameSyncPack {

	boxs := []*typings.TGameBox{}

	for _, line := range game.GameMap {
		for _, block := range line {
			if block.Box != nil && !block.Box.Hasdestoryed {
				boxs = append(boxs, block.Box)
			}
		}
	}

	return &TGameSyncPack{
		Props:   boxs,
		Players: game.Players,
		Bubbles: game.Bubbles,
	}
}

func (game *Game) getPlayerByName(name string) *typings.TGamePlayer {
	for _, v := range game.Players {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func CreateGame(room *typings.Room) {
	// 生成地图数据
	boxs, roles, mymap := gamemap.GetGameMap()
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
			Role:       p.Role,
		})

		i++
	}

	gameInfo := &Game{
		RoomId:  room.Id,
		Props:   boxs,
		Players: players,
		Bubbles: []*typings.TGameBubble{},
		GameMap: mymap,
	}

	gameMap[room.Id] = gameInfo
}
