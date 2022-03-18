package main

import (
	"github.com/NYTimes/gizmo/server/kit"
	"os"
	gmo "using-kit/gizmo"
)

//bubble up comments for clarity
// Run will use environment variables to configure the server then register the given
// Service and start up the server(s).
// This will block until the server shuts down.
func main() {
	os.Setenv("PORT", "8833")
	os.Setenv("GIZMO_SKIP_OBSERVE", "true")

	err := kit.Run(gmo.NewService())
	if err != nil {
		panic("problems running service: " + err.Error())
	}
}
