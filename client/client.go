package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"idea_server/protolcol/message"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	//u := url.URL{Scheme: "wss", Host: "idea.vinf.top:8001", Path: "/im"}
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8001", Path: "/auth/ws"}
	log.Printf("connecting to %s", u.String())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQ2NDEwOTAsImlkIjoiMSIsIm9yaWdfaWF0IjoxNjM2MDAxMDkwfQ.UXY0VzBpFMhLU9M-2Rt5KpAVB6_eb3kLyG1jLMh3ZEk"
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, bytes, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			msg := &message.Msg{}
			err = proto.Unmarshal(bytes, msg)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(msg)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			msg := &message.Msg{
				Mid:     0,
				Uid:     1,
				IsGroup: true,
				Did:     1,
				Op:      message.SendMsg,
				Data:    nil,
			}
			bytes, err := msg.Marshal()
			if err != nil {
				log.Println(err)
				return
			}
			err = c.WriteMessage(websocket.BinaryMessage, bytes)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
