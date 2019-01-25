package main

import (
	"fmt"
)

type Message interface {
	Content() []byte
	// Can use visitor pattern here.
	// SMSContent() []byte
	// BrowserNotificationContent() []byte
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
	title string
	body  string
}

func (p *PromotionMessage) Content() []byte {
	msg := fmt.Sprintf("%s - %s", p.title, p.body)
	return []byte(msg)
}

func main() {
	msgs := []Message{
		&GreetingMessage{[]byte("hello world")},
		&PromotionMessage{title: "CNY Sale", body: "50% Off"},
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
