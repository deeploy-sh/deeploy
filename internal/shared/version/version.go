package version

// Version is set via -ldflags during release builds.
// In local development it defaults to "dev".
//
// Build example:
//
//	go build -ldflags "-X github.com/deeploy-sh/deeploy/internal/shared/version.Version=v0.1.0"
var Version = "dev"
