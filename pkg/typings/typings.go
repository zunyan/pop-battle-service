package typings

import "time"

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
	Gridx        int           `json:"gridX"`
	Gridy        int           `json:"gridY"`
	Status       bool          `json:"status"` // 是否已经被炸开
	Props        TGamePropEnum `json:"props"`  // 道具的枚举属性
	Hasdestoryed bool
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
	Bubbles    int               `json:"bubbles"`    // 最大可放置的泡泡
	MoveTarget string            `json:"moveTarget"` // 移动方向 Left, Right, Up, Down
}

type TGameBubble struct {
	Gridx        int           `json:"gridX"`      // 泡泡X坐标
	Gridy        int           `json:"gridY"`      // 泡泡y坐标
	Power        int           `json:"power"`      // 泡泡威力
	CreateTime   time.Duration `json:"createtime"` // 创建时间
	Hasdestoryed bool          // 是否被销毁
}

type TGameBoomBubble struct {
	Gridx  int `json:"gridX"`
	Gridy  int `json:"gridY"`
	Left   int `json:"left"`
	Right  int `json:"right"`
	Top    int `json:"top"`
	Bottom int `json:"bottom"`
}
type TGameInfo struct {
	Props   []*TGameBox    `json:"props"`
	Players []*TGamePlayer `json:"players"`
	Bubbles []*TGameBubble `json:"bubbles"`
}
