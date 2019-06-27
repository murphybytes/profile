package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/murphybytes/profiler"
	"github.com/murphybytes/profiler/config"
)

func main() {
	fmt.Printf("pid %d\n", os.Getpid())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	settings, err := config.New()
	if err != nil {
		log.Fatal("could not create new config", err)
	}
	profiler.Heap(ctx, settings)

	<-ch
	cancel()

	fmt.Println("exiting")
}
