package messenger

import "context"

type Messenger interface {
	Name() string
	Ping(ctx context.Context) error
	Send(ctx context.Context, message Message) error
}

type Message struct {
	Body string
}
