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
