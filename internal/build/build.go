// Package build tells what build of tuna runs: its version, and whether
// it is the server build (cmd/tuna-server), which runs tunad and the CLI
// alone, with no app, and updates itself.
package build

import "github.com/egoist/mygo"

// Server is set by the server build.
var Server bool

// Version is the version of tuna that runs: the app's, or the one the
// server build was made of.
func Version() string { return mygo.App.Version() }

// Daemon is what tunad's code is made of, which the app sets (tunadBuild)
// and tunad tells clients (proto.Welcome.Build): an update that leaves it
// as it was needs no restart of tunad. The server build leaves it empty.
var Daemon string
