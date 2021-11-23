package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		//升级get请求为webSocket协议
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		//返回前关闭
		defer ws.Close()

		// 开始监听websocket 消息
		for {
			//读取ws中的数据
			_, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			var msg interface{}
			fmt.Println(string(message))
			jsonErr := json.Unmarshal(message, &msg)
			if jsonErr != nil {
				fmt.Println(jsonErr)
			} else {
				fmt.Println(msg)
			}
			ws.WriteMessage(1, []byte("hello"))
		}
	})
	r.Run(fmt.Sprintf(":%d", 9502))
}
