package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/72636c/stratus/internal/config"
)

// func main() {
// 	app, err := cli.New()
// 	check(err)

// 	if app == nil {
// 		return
// 	}

// 	ctx := context.Background()

// 	err = app.Do(ctx)
// 	check(err)
// }

// func check(err error) {
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// }

func main() {
	err := Do()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Do() error {
	file, err := ioutil.ReadFile("internal/config/testdata/template.json")
	if err != nil {
		return err
	}

	template, err := config.TraverseNestedStacks(file)
	if err != nil {
		return err
	}

	fmt.Println(string(template))

	return nil
}
