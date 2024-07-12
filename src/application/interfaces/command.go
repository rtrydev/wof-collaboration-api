package interfaces

import "context"

type CommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}
