package cli

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/72636c/stratus/internal/context"
)

func Do(
	ctx context.Context,
	args *Args,
) error {
	for index := 0; index < len(args.cfg.Stacks); index++ {
		err := do(ctx, args, index)
		if err != nil {
			return err
		}
	}

	return nil
}

func do(
	ctx context.Context,
	args *Args,
	index int,
) error {
	stack := args.cfg.Stacks[index]

	group, ctx := errgroup.WithContext(ctx)

	output := make(chan interface{})

	group.Go(func() (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = fmt.Errorf("recovered from panic: %+v", recovered)
			}
		}()

		defer close(output)

		output <- fmt.Sprintf("Stratus.[%d].StackConfig", index)

		output <- stack

		ctx := context.WithOutput(ctx, output)

		return args.command(ctx, args.client, stack)
	})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	group.Go(func() error {
		for {
			select {
			case <-ticker.C:
				fmt.Printf(".")

			case message, ok := <-output:
				args.logger(message, ok)

				if !ok {
					return nil
				}
			}
		}
	})

	return group.Wait()
}
