package dataflash

const (
	// Head1 is
	Head1 byte = 0xA3
	// Head2 is
	Head2 byte = 0x95
)

// Message is
type Message struct {
}

// MessageHeader is
type MessageHeader struct {
	MessageID byte
}

func NewMessage() *Message {
	return &Message{}
}
