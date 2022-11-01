package schema

type ChatMessage struct {
	ID         uint
	Content    string
	ReceiverID uint
	SenderID   uint
}
