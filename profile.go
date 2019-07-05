// Package profiler contains function that will perform pprof profile operations.
package profile

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/murphybytes/profile/config"
)

// Mutex stack traces of holders of contended mutexes
func Mutex(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.MutexProfileName)
	profile(ctx, "mutex", outputFile, settings.MutexProfilerSignal)
}

// Block stack traces that led to blocking on synchronization primitives
func Block(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.BlockProfileName)
	profile(ctx, "block", outputFile, settings.BlockProfilerSignal)
}

// ThreadCreate outputs stack traces that led to the creation of new OS threads
func ThreadCreate(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.ThreadCreateProfileName)
	profile(ctx, "threadcreate", outputFile, settings.ThreadCreateProfilerSignal)
}

// Allocs output a profile containing memory allocations
func Allocs(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.AllocsProfileName)
	profile(ctx, "allocs", outputFile, settings.AllocsProfilerSignal)
}

// Goroutine outputs a goroutine profile.
func Goroutine(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.GoroutineProfileName)
	profile(ctx, "goroutine", outputFile, settings.GoroutineProfilerSignal)
}

// Heap outputs a heap profile.
func Heap(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.HeapProfileName)
	profile(ctx, "heap", outputFile, settings.HeapProfilerSignal)
}

func Run(profileFunc func(ctx context.Context, cfg *config.Settings) ) {
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
			cancel()
		}()

		cfg, err := config.New()
		if err != nil {
			log.Printf("unable to fetch profiler configuration %q", err)
		}
		profileFunc(ctx, cfg)
		<-ch
	}()
}

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
