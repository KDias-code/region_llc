package tasks

import "context"

type Cache interface {
	Status(ctx context.Context, id string, data Entity) (err error)
}
