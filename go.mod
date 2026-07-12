module github.com/magnobit/quell-cli

go 1.25

require (
	github.com/magnobit/quell v0.0.4
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Points at the local quell working copy so this builds against the new
// public execute/compile APIs (Extra passthrough, AWS/Google credential
// fixes, execute package itself) before they're tagged in a real quell
// release. Must be removed — replaced with a real `require
// github.com/magnobit/quell vX.Y.Z` once quell publishes a version with
// these changes — before this repo is pushed/published anywhere public;
// a public consumer has no access to this local path.
replace github.com/magnobit/quell => ../quell
