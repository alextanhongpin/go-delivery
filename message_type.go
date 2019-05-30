package omni

type MessageType string

var (
	Browser MessageType = "browser"
	Chatbot MessageType = "chatbot"
	Email   MessageType = "email"
	Mobile  MessageType = "mobile"
	SMS     MessageType = "sms"
)
