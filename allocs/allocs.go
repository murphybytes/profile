// Package allocs use underscore import to add allocs profiling to your application
package allocs

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.Allocs)
}