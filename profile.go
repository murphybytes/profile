// Package profiler contains function that will perform pprof profile operations.
package profile

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/murphybytes/profile/config"
)

func profile(ctx context.Context, profileName, fileName string, sig config.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.Signal(sig))

	go func() {
		defer func() {
			signal.Stop(ch)
			close(ch)
		}()

		for {
			select {
			case <-ch:
				if err := generate(profileName, fileName); err != nil {
					log.Println(profileName, "profile failed with", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func generate(profileName, fileName string) error {
	if profileName == "heap" {
		runtime.GC()
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	prof := pprof.Lookup(profileName)
	if prof == nil {
		return fmt.Errorf("could not create %q profile, no such profiler exists", profileName)
	}
	if err = prof.WriteTo(f, 0); err != nil {
		return err
	}
	return nil
}
