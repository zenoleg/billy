package transport

import "context"

type Bot interface {
	Run(ctx context.Context) error
}
