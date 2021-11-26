package gamemap

import (
	"encoding/json"
	"fmt"
	"os"
	"pop-battle-service/pkg/typings"
)

type TGameMapBlock struct {
	Floor string `json:"floor"`
	Top   string `json:"top"`
	Type  bool   `json:"candestory"`
	Prop  string `json:"prop"`

	CanDestory bool
	Bubble     *typings.TGameBubble
	Box        *typings.TGameBox
}

func GetGameMap() ([]*typings.TGameBox, [][]int, [][]*TGameMapBlock) {

	file, fileError := os.Open("./priate.json") // For read access.

	if fileError != nil {
		fmt.Println(fileError)
		return nil, nil, nil
	}

	stat, _ := file.Stat()

	var gamemap [][]*TGameMapBlock
	buf := make([]byte, stat.Size())

	file.Read(buf)
	json.Unmarshal(buf, &gamemap)

	// if jsonErr != nil {
	// }

	boxs := []*typings.TGameBox{}
	roles := [][]int{}
	for y, line := range gamemap {
		for x, block := range line {
			block.Bubble = nil
			if block.Top == "role" {
				roles = append(roles, []int{x, y})
			} else if block.Top != "" {
				box := &typings.TGameBox{
					Gridx:  x,
					Gridy:  y,
					Status: false,
					Props:  typings.TGamePropEnum_BUBBLES,
				}
				boxs = append(boxs, box)
				block.Box = box
			}
		}
	}

	return boxs, roles, gamemap
}
