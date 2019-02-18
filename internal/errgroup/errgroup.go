package errgroup

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	*errgroup.Group
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	group, ctx := errgroup.WithContext(ctx)
	return &Group{Group: group}, ctx
}

func (group *Group) Go(function func() error) {
	group.Group.Go(func() (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = fmt.Errorf("recovered from panic: %+v", recovered)
			}
		}()

		return function()
	})
}
