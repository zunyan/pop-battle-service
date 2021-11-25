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
	Type  bool   `json:"type"`
	Prop  string `json:"prop"`
}

func GetGameMap() ([]*typings.TGameBox, [][]int) {

	file, fileError := os.Open("./priate.json") // For read access.

	if fileError != nil {
		fmt.Println(fileError)
		return nil, nil
	}

	stat, _ := file.Stat()

	var gamemap [][]TGameMapBlock
	buf := make([]byte, stat.Size())

	file.Read(buf)
	json.Unmarshal(buf, &gamemap)

	// if jsonErr != nil {
	// }

	boxs := []*typings.TGameBox{}
	roles := [][]int{}
	for y, line := range gamemap {
		for x, block := range line {
			if block.Top == "role" {
				roles = append(roles, []int{x, y})
			}
			if block.Top != "" {
				boxs = append(boxs, &typings.TGameBox{
					Gridx:  x,
					Gridy:  y,
					Status: false,
					Props:  typings.TGamePropEnum_BUBBLES,
				})
			}
		}
	}

	return boxs, roles
}
