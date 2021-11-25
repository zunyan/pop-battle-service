package main

import (
	"fmt"
	"log"
	"net/http"
	"pop-battle-service/pkg/game"
	"pop-battle-service/pkg/hall"
	"pop-battle-service/pkg/room"
	"pop-battle-service/pkg/store"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

type SocketACK bool

func main() {

	store.Init()

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	// 默认连接, 不可删除
	server.OnConnect("/", func(s socketio.Conn) error { return nil })
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {})
	server.OnError("/", func(s socketio.Conn, e error) {
		s.Close()
		fmt.Println("meet error:", e)
	})

	hall.LinkRouter(server) //大厅
	room.LinkRouter(server) // 房间服务
	game.Init()(server)     // 对局服务

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000 ...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
