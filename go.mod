module github.com/magnobit/quell-cli

go 1.25

require (
	github.com/magnobit/quell v0.0.6
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// TEMPORARY, local-only: points at the local quell working copy so this
// builds against the new public `simulate` package (local circuit
// simulator, backing `quell simulate` / `quell run --backend local` below)
// before it's tagged in a real quell release. v0.0.6 (the require line
// above) already has everything else this repo needs — this replace exists
// solely for the new simulate package. Must be removed — swap back to a
// real `require github.com/magnobit/quell vX.Y.Z` once quell tags a version
// that includes simulate/ — before this repo is pushed/published; a public
// consumer has no access to this local path. Same two-step dance as the
// v0.0.4 -> v0.0.6 upgrade: tag+push quell first, then point this back at
// the real version and delete this replace line.
replace github.com/magnobit/quell => ../quell
