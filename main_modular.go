// +build modular

package main

import (

	// NOTE: If compiling gopherbot yourself, you can comment out or remove
	// most of the import lines below to shrink the binary or remove unwanted
	// or unneeded funcationality. You'll need at least one connector for your
	// bot to be useful, however.

	// Many included plugins already have 'Disabled: true', but you can also
	// disable by adding that line to conf/plugins/<plugname>.yaml

	// *** Included Authorizer plugins
	_ "github.com/lnxjedi/gopherbot/goplugins/groups"

	// *** Included Go plugins, of varying quality
	_ "github.com/lnxjedi/gopherbot/goplugins/help"
	_ "github.com/lnxjedi/gopherbot/goplugins/links"
	_ "github.com/lnxjedi/gopherbot/goplugins/lists"
	_ "github.com/lnxjedi/gopherbot/goplugins/ping"
	/* Enable profiling. This inflates the binary when enabled, but if the
	   robot ever stops responding for any reason, it's handy for getting a
	   dump of all goroutines. Example usage:

	   $ go tool pprof http://localhost:8888/debug/pprof/goroutine
	   ...
	   Entering interactive mode (type "help" for commands, "o" for options)
	   (pprof) list lnxjedi
	   Total: 11
	   ROUTINE ======================== github.com/lnxjedi/gopherbot/bot...
	*/// _ "net/http/pprof"
)
