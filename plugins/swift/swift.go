// Package swift implements the "swift" runtime.
package swift

import (
	"github.com/apex/apex/function"
	"github.com/apex/apex/plugins/nodejs"
)

func init() {
	function.RegisterPlugin("swift", &Plugin{})
}

const (
	// Runtime name used by Apex
	Runtime = "swift"
)

// Plugin implementation.
type Plugin struct{}

// Open adds the shim and swift defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = "swift build -c release; mv .build/release/main main; mv .build/release/*.so .; rm -rf .build; rm -rf Packages"
	}

	fn.Shim = true
	fn.Runtime = nodejs.Runtime

	if fn.Hooks.Clean == "" {
		fn.Hooks.Clean = "swift build --clean=dist"
	}

	return nil
}
