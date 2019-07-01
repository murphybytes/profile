// Package threadcreate use underscore import to create a profile of stack traces that led to the creation of new OS threads.
package threadcreate

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.ThreadCreate)
}