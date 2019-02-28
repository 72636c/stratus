package main

import (
	"context"
	"fmt"
)

func test() {
	ctx, _ := context.WithCancel(context.Background())

	fmt.Printf("%T\b", ctx)
}
