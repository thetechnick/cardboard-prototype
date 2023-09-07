package steps

import (
	"context"
	"io"
)

type Interface interface {
	Name() string
	Run(ctx context.Context, out io.Writer) error
}

type List []Interface

func (l List) Run(ctx context.Context, out io.Writer) error {
	for _, t := range l {
		if err := t.Run(ctx, out); err != nil {
			return err
		}
	}
	return nil
}
