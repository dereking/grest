package controllers

import (
	"time"

	"github.com/dereking/grest/mvc"

	"golang.org/x/net/websocket"
)

type WsController struct {
	mvc.Controller
}

func (c *WsController) Chat(ws *websocket.Conn) {

	defer ws.Close()

	var err error
	var str string

	for {
		/*if err = websocket.Message.Receive(ws, &str); err != nil {
			break
		} else {
			fmt.Println("从客户端收到：", str)
		}*/

		str = "hello, I'm server."

		if err = websocket.Message.Send(ws, str); err != nil {
			break
		} else {
			time.Sleep(time.Second * 2)
		}
	}
}
