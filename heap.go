package profiler

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/murphybytes/profiler/config"
)

// Heap outputs a heap profile to a file given a signal. Heap expects a cancel context as an argument. The default
// file is heap.profile in the current working directory. The triggering signal is SIGRTMAX-14 (50). These can both
// be changed. See Settings.
func Heap(ctx context.Context, settings config.Settings) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, settings.HeapProfilerSignal)

	go func() {
		defer func() {
			close(ch)
		}()

		for {
			select {
			case <-ch:
				memProfileFileName := filepath.Join(settings.ProfileDirectory, settings.HeapProfileFileName)
				if err := profileHeap(memProfileFileName); err != nil {
					log.Print("heap profile failed with", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

}

func profileHeap(memProfileFileName string) error {
	runtime.GC()
	f, err := os.Create(memProfileFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	prof := pprof.NewProfile("heap")
	if prof == nil {
		return errors.New("could not create heap profile, no such profiler exists")
	}
	if err = prof.WriteTo(f, 0); err != nil {
		return err
	}
	return nil
}