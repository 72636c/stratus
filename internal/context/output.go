package context

import (
	"log"
	"sync"
)

var (
	defaultOutput     = make(chan interface{})
	defaultOutputOnce = new(sync.Once)
)

func newOutput() chan interface{} {
	defaultOutputOnce.Do(func() {
		go func() {
			for message := range defaultOutput {
				log.Println(message)
			}
		}()
	})

	return defaultOutput
}

func Output(ctx Context) chan<- interface{} {
	output, ok := ctx.Value(outputKey).(chan<- interface{})
	if ok {
		return output
	}

	return newOutput()
}

func WithOutput(ctx Context, output chan<- interface{}) Context {
	return WithValue(ctx, outputKey, output)
}
