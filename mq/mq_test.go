package mq

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestManager_Consume(t *testing.T) {
	m := &Manager{}
	err := m.Init()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go m.Consume(ctx, func(bytes []byte) {
		log.Println("1:", string(bytes))
		time.Sleep(time.Second)
	})

	for i := 0; i < 10; i++ {
		m.Product([]byte(strconv.FormatInt(int64(i), 10)))
	}
	time.Sleep(12 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)

	m.Close()
}
