package handlers

import (
	"fmt"

	"github.com/danielmadu/goexpose/config"
	"golang.org/x/net/websocket"
)

func WebSocket(ws *websocket.Conn) {
	defer ws.Close()
	channel := config.GetChannel()
	for {
		message := <-channel
		encoded, err := message.Encode()
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err := ws.Write(encoded); err != nil {
			fmt.Println(err)
			ws.Close()
			return
		}

		var buff = make([]byte, 1048576)
		messageResponse := config.Message{}
		n, err := ws.Read(buff)

		if err != nil {
			continue
		}

		err = messageResponse.Decode(buff[:n])

		if err != nil {
			fmt.Println(err)
			continue
		}

		channel <- messageResponse
	}
}
