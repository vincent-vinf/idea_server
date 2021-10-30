package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	//u := url.URL{Scheme: "wss", Host: "idea.vinf.top:8001", Path: "/connect"}
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8001", Path: "/auth/connect"}
	log.Printf("connecting to %s", u.String())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzU1NzY3MDQsImlkIjoiMSIsIm9yaWdfaWF0IjoxNjM1NTczMTA0fQ.VcNv7g_mvpe2NEckJ-ci4INSlbpeY-BQn1ZZEIVkNfQ"
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
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte("ping"))
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
