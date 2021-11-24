package store

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	ROOM_STATUS_WAITING = 0
	ROOM_STATUS_IN_GAME = 1

	// 用户状态
	PLAYER_STATUS_PENDING = 0
	PLAYER_STATUS_READY   = 1
)
const (
	ROLE_MARID = "role_marid"
	ROLE_BUZZI = "role_buzzi"
	ROLE_DAO   = "role_dao"
	ROLE_CAPPI = "role_cappi"
)

type Player struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Status   int    `json:"status"`
	IsMaster bool   `json:"isMaster"`
}
type Room struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Players     []*Player `json:"players"`
	TotalPlayer int       `json:"totalPlayer"`
	Status      int       `json:"status"`
}

var RoomList []*Room
var RoomListMap map[string]*Room

func Init() {
	RoomList = []*Room{}
	RoomListMap = map[string]*Room{}
}

func CreateRoom(roomName string, username string) *Room {
	room := &Room{
		Id:   uuid.New().String(),
		Name: roomName,
		Players: []*Player{
			&Player{Name: username, IsMaster: true, Status: PLAYER_STATUS_PENDING, Role: ROLE_MARID},
		},
		TotalPlayer: 4,
		Status:      ROOM_STATUS_WAITING,
	}
	RoomList = append(RoomList, room)
	RoomListMap[room.Id] = room

	return room
}

func HasRoom(roomId string) bool {
	_, exist := RoomListMap[roomId]
	return exist
}

func FindPlayerInRoom(room *Room, username string) *Player {
	var player *Player
	for _, v := range room.Players {
		if v.Name == username {
			player = v
			break
		}
	}

	return player
}

func JoinRoom(roomId string, username string) (*Room, error) {
	if !HasRoom(roomId) {
		return nil, errors.New("无效的房间号")
	}

	room := RoomListMap[roomId]

	if room.TotalPlayer <= len(room.Players) {
		return nil, errors.New("房间已经满人")
	}

	if FindPlayerInRoom(room, username) != nil {
		fmt.Println("玩家已经存在与该房间", room.Name, username)
		return room, nil
	}

	player := &Player{
		Name:   username,
		Status: PLAYER_STATUS_PENDING,
		Role:   ROLE_MARID,
	}

	room.Players = append(room.Players, player)
	return room, nil
}
