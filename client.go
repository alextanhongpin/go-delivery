package main

import (
	"fmt"
)

type Message interface {
	Content() []byte
	// From, To, Time
}

type Client interface {
	Send(Message) error
}

type SMSClient struct{}

func (s *SMSClient) Send(msg Message) error {
	fmt.Println("[SMS]", string(msg.Content()))
	return nil
}

type BrowserNotificationClient struct{}

func (b *BrowserNotificationClient) Send(msg Message) error {
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

	clients := []Client{
		new(BrowserNotificationClient),
		new(SMSClient),
	}

	for _, client := range clients {
		for _, msg := range msgs {
			client.Send(msg)
		}
	}
}
