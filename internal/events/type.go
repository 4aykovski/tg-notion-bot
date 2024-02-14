package events

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Data interface{}
	Meta interface{}
}
