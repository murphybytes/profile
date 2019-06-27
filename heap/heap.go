// Package heap provides an easy way to add heap profiling to your application
package heap

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/murphybytes/profiler"
	"github.com/murphybytes/profiler/config"
)

func init() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
		defer func() {
			signal.Stop(ch)
			close(ch)
		}()
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-ch
			fmt.Println("exiting heap handler")
			cancel()
		}()

		cfg, err := config.New()
		if err != nil {
			log.Println("unable to fetch heap profiler configuration", err)
		}
		profiler.Heap(ctx, cfg)
		<-ch
	}()
}
