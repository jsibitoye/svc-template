package version

var (
	// Overridden at build time via -ldflags.
	Version   = "0.1.0"
	GitCommit = "dev"
	BuildDate = "unknown"
)