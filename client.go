package main

import (
	"fmt"
)

type Message interface {
	Content() []byte
}

type Notification interface {
	Send(Message) error
}

type BrowserNotification struct {}

func (b *BrowserNotification) Send(msg Message) error {
	fmt.Println("[BrowserNotification]", string(msg.Content()))
	return nil
}

type GreetingMessage struct {
	msg []byte
}

func (g *GreetingMessage) Content() []byte {
	return g.msg
}

type PromotionMessage struct {
	msg []byte
}

func (p *PromotionMessage) Content() []byte {
	return p.msg
}

func main() {
	msgs := []Message{
		&GreetingMessage{[]byte("hello world")},
		&PromotionMessage{[]byte("50% off")},
	}

	client := new(BrowserNotification)
	for _, msg := range msgs {
		client.Send(msg)
	}
}
