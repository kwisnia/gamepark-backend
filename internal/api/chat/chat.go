package chat

type Message struct {
	ID         uint
	Content    string
	ReceiverID uint
	SenderID   uint
}
