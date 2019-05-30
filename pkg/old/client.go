package main

import (
	"fmt"
)

type Message interface {
	Content() []byte

	// Can use visitor pattern here.
	SMSContent() []byte
	BrowserNotificationContent() []byte

	// From, To, Time
}

type Client interface {
	Send(Message) error
}

type SMSClient struct{}

func (s *SMSClient) Send(msg Message) error {
	fmt.Printf("[SMS] %s\n", msg.SMSContent())
	return nil
}

type BrowserNotificationClient struct{}

func (b *BrowserNotificationClient) Send(msg Message) error {
	fmt.Printf("[BrowserNotification] %s\n", msg.BrowserNotificationContent())
	return nil
}

type GreetingMessage struct {
	msg string
}

func (g *GreetingMessage) Content() []byte {
	return []byte(g.msg)
}

func (g *GreetingMessage) SMSContent() []byte {
	return g.Content()
}

func (g *GreetingMessage) BrowserNotificationContent() []byte {
	msg := fmt.Sprintf("<h1>%s</h1>", g.msg)
	return []byte(msg)
}

type PromotionMessage struct {
	title string
	body  string
}

func (p *PromotionMessage) Content() []byte {
	msg := fmt.Sprintf("%s - %s", p.title, p.body)
	return []byte(msg)
}

func (p *PromotionMessage) SMSContent() []byte {
	return p.Content()
}

func (p *PromotionMessage) BrowserNotificationContent() []byte {
	msg := fmt.Sprintf("<h1>%s</h1><p>%s</p>", p.title, p.body)
	return []byte(msg)
}

func main() {
	msgs := []Message{
		&GreetingMessage{"hello world"},
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
