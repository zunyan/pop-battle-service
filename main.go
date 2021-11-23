package main

import (
	"wsService/pkg/websocket"
)

func Init() {
	//启动ws 监听c端传过来的消息
	websocket.Start()
}

func main() {
	//启动
	Init()
	// time.Sleep(time.Minute)
	//game.Test()
}
