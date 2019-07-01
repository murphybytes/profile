package main

import (
	"fmt"
	"os"
	"os/signal"

	_ "github.com/murphybytes/profile/allocs"
	_ "github.com/murphybytes/profile/goroutine"
	_ "github.com/murphybytes/profile/heap"
)

func main() {
	fmt.Printf("pid %d\n", os.Getpid())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("exiting")
}
