// Package cpu facilitates the use of underscore imports to include pprof CPU profiling in an application.
// various parameters can be tweaked to customize the profile.  See config.
package cpu

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.CPU)
}
