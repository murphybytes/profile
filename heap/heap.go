// Package heap provides an easy way to add heap profiling to your application
package heap

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.Heap)
}
