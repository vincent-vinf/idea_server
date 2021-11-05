package mq

import (
	"context"
	"github.com/streadway/amqp"
	"idea_server/util"
	"log"
	"sync"
)

const (
	queueName = "idea"
)

var (
	manager *Manager
	once    sync.Once
)

func GetInstance() *Manager {
	if manager == nil {
		once.Do(func() {
			manager = &Manager{}
			err := manager.Init()
			if err != nil {
				panic(err)
			}
		})
	}
	return manager
}

type Manager struct {
	conn      *amqp.Connection
	sendCh    *amqp.Channel
	receiveCh *amqp.Channel
}

func (m *Manager) Init() error {
	conn, err := amqp.Dial(util.LoadMqCfg().Url)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		return err
	}
	m.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}
	m.sendCh = ch

	ch, err = conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}
	m.receiveCh = ch
	err = m.receiveCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}
	return nil
}

func (m *Manager) Close() {
	m.conn.Close()
	m.sendCh.Close()
}

func (m *Manager) Product(bytes []byte) error {
	err := m.sendCh.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         bytes,
		})
	return err
}

func (m *Manager) Consume(ctx context.Context, process func([]byte)) {
	msgQueue, err := m.receiveCh.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case d := <-msgQueue:
			process(d.Body)
		case <-ctx.Done():
			return
		}
	}
}
