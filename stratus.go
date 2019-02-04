package main

import (
	"fmt"
	"os"

	"github.com/72636c/stratus/internal/cli"
	"github.com/72636c/stratus/internal/context"
)

func main() {
	app, err := cli.New()
	check(err)

	if app == nil {
		return
	}

	ctx := context.Background()

	err = app.Do(ctx)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
