// Package block use underscore import to add a profile of stack traces that led to blocking on synchronization
// primitives
package block

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.Block)
}
