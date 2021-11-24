package store

import (
	"github.com/google/uuid"
)

const (
	ROOM_STATUS_WAITING = 0
	ROOM_STATUS_IN_GAME = 1

	// 用户状态
	PLAYER_STATUS_PENDING = 0
	PLAYER_STATUS_READY   = 1
)

type Player struct {
	Name      string `json:"name"`
	RoleIndex int    `json:"roleIndex"`
	Status    int    `json:"status"`
	IsMaster  int    `json:"isMaster"`
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
			&Player{Name: username, IsMaster: true, Status: PLAYER_STATUS_PENDING},
		},
		TotalPlayer: 4,
		Status:      ROOM_STATUS_WAITING,
	}
	RoomList = append(RoomList, room)
	RoomListMap[room.Id] = room

	return room
}
