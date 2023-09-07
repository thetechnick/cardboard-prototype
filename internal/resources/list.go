package resources

import (
	"context"

	"cardboard.package-operator.run/internal"
)

type List []Interface

func (l List) Watch(ctx context.Context) (<-chan WatchEvent, error) {
	chs := make([]<-chan WatchEvent, len(l))
	for i := range l {
		var err error
		chs[i], err = l[i].Watch(ctx)
		if err != nil {
			return nil, err
		}
	}
	return internal.MergeCh(chs...), nil
}
