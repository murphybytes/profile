// Package trace will generate output for the Go execution tracer if underscore imported into an application.
package trace

import "github.com/murphybytes/profile"

func init() {
	profile.Run(profile.Trace)
}
