// Package goroutine is used provide goroutine profiling using an underscore import.
package goroutine

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.Goroutine)
}
