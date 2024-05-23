package messenger

import "context"

type Messenger interface {
	Name() string
	Send(ctx context.Context, message Message) error
}

type Message struct {
	Body string
}
