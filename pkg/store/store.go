package store

import (
	"errors"
	"fmt"
	"pop-battle-service/pkg/typings"

	"github.com/google/uuid"
)

var RoomList []*typings.Room
var RoomListMap map[string]*typings.Room

func Init() {
	RoomList = []*typings.Room{}
	RoomListMap = map[string]*typings.Room{}
}

func CreateRoom(roomName string, username string) *typings.Room {
	player := &typings.Player{
		Name:     username,
		IsMaster: true,
		Status:   typings.PLAYER_STATUS_PENDING,
		Role:     typings.ROLE_MARID,
	}
	room := &typings.Room{
		Id:          uuid.New().String(),
		Name:        roomName,
		Players:     []*typings.Player{player},
		TotalPlayer: 4,
		Status:      typings.ROOM_STATUS_WAITING,
	}
	RoomList = append(RoomList, room)
	RoomListMap[room.Id] = room

	return room
}

func HasRoom(roomId string) bool {
	_, exist := RoomListMap[roomId]
	return exist
}

func FindPlayerInRoom(room *typings.Room, username string) *typings.Player {
	var player *typings.Player
	for _, v := range room.Players {
		if v.Name == username {
			player = v
			break
		}
	}

	return player
}

func JoinRoom(roomId string, username string) (*typings.Room, error) {
	if !HasRoom(roomId) {
		return nil, errors.New("无效的房间号:" + roomId)
	}

	room := RoomListMap[roomId]

	if room.TotalPlayer <= len(room.Players) {
		return nil, errors.New(room.Name + "房间已经满人")
	}

	if FindPlayerInRoom(room, username) != nil {
		fmt.Printf("玩家%v已经存在于%v房间\n", username, room.Name)
		return room, nil
	}

	if room.Status == typings.ROOM_STATUS_IN_GAME {
		return nil, errors.New("游戏进行中，无法加入")
	}

	player := &typings.Player{
		Name:   username,
		Status: typings.PLAYER_STATUS_PENDING,
		Role:   typings.ROLE_MARID,
	}

	room.Players = append(room.Players, player)
	return room, nil
}
