// tuna-server is the server build of tuna: tunad and the tuna CLI,
// without the app, for a machine whose sessions tuna on another one shows
// (Settings › Hosts), as a Linux box with no display. It installs in
// ~/.local with install-server.sh, which runs tunad as a service, and
// updates itself (tuna daemon update, or the app's Update tuna).
//
// Its version comes with the build:
//
//	go build -ldflags "-X main.version=2.1.0 -X github.com/egoist/mygo.production=1" -o tuna ./cmd/tuna-server
package main

import (
	"fmt"
	"os"

	"github.com/egoist/mygo"

	"tuna/internal/build"
	"tuna/internal/cli"
	"tuna/internal/daemon"
)

// version is the version of tuna this build is of, set with -ldflags.
var version = "dev"

func main() {
	build.Server = true
	mygo.App.SetVersion(version)
	if len(os.Args) > 1 && os.Args[1] == "__tunad" {
		if err := daemon.Run(daemon.Options{}); err != nil {
			fmt.Fprintln(os.Stderr, "tunad:", err)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) < 2 {
		os.Exit(cli.Run([]string{"help"}))
	}
	if !cli.IsCommand(os.Args[1]) {
		fmt.Fprintf(os.Stderr, "tuna: unknown command %q (see tuna help)\n", os.Args[1])
		os.Exit(2)
	}
	os.Exit(cli.Run(os.Args[1:]))
}
