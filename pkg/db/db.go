package db

import "context"

// See here:
type Database interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error

	Query(ctx context.Context, query Query) Rows
	Execute(ctx context.Context, query Query) Result
}
