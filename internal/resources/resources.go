// Resources are inputs or outputs to Jobs.
// A resource generally supports:
// - Getting the resource
// - Update or create the resource
// - Watch the resource for changes
package resources

import (
	"context"
)

type Interface interface {
	Name() string
	Get(context.Context) (Volume, error)
	Put(context.Context, Volume) error
	Watch(context.Context) (<-chan WatchEvent, error)
}

type Volume struct {
	Path string
}

type WatchEvent struct {
	Err error
}
