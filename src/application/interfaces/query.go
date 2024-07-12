package interfaces

import "context"

type QueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}
