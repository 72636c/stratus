package main

import (
	"fmt"
	"os"

	"github.com/72636c/stratus/internal/cli"
	"github.com/72636c/stratus/internal/context"
)

func main() {
	args, err := cli.FromCommandLine()
	check(err)

	ctx := context.Background()

	err = cli.Do(ctx, args)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
