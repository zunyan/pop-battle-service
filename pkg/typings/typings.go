package typings

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

// TGamePropEnum
const (
	TGamePropEnum_NONE    = 0
	TGamePropEnum_SHOSE   = 1
	TGamePropEnum_LOTION  = 2
	TGamePropEnum_BUBBLES = 3
)

type TGamePropEnum uint

type TGameBox struct {
	Gridx  int           `json:"gridX"`
	Gridy  int           `json:"gridY"`
	Status bool          `json:"status"`
	Props  TGamePropEnum `json:"props"`
}

// TGamePlayerStatus
type TGamePlayerStatus = uint

const (
	TGamePlayerStatus_ALIVE = 1
	TGamePlayerStatus_DEAD  = 2
)
const (
	TGamePlayerMoveTarget_Left  = "Left"
	TGamePlayerMoveTarget_Right = "Right"
	TGamePlayerMoveTarget_Up    = "Up"
	TGamePlayerMoveTarget_Down  = "Down"
	TGamePlayerMoveTarget_None  = "None"
)

type TGamePlayer struct {
	X          int               `json:"x"`
	Y          int               `json:"y"`
	Gridx      int               `json:"gridX"`
	Gridy      int               `json:"gridY"`
	Name       string            `json:"name"`
	Status     TGamePlayerStatus `json:"status"`
	Speed      int               `json:"speed"`
	Power      int               `json:"power"`
	Bubbles    int               `json:"bubbles"`
	MoveTarget string            `json:"moveTarget"`
}

type TGameBubble struct {
	Gridx int `json:"gridX"`
	Gridy int `json:"gridY"`
	Power int `json:"power"`
}
type TGameInfo struct {
	Props   []*TGameBox    `json:"props"`
	Players []*TGamePlayer `json:"players"`
	Bubbles []*TGameBubble `json:"bubbles"`
}
