// Package mutex underscore import the package to add profiling of contended mutexes.
package mutex

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.Mutex)
}