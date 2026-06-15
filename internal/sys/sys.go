package sys

// Version is the version information
var Version = "0.0.1"

// Build is the hash of build
var Build = "2ed26fe1"

// OsArch is the architecture of the build
var OsArch = "linux/amd64"

// GetVersion returns the application version
func GetVersion() string {
	return Version
}

// GetBuild returns the build hash
func GetBuild() string {
	return Build
}

// GetOsArch returns the build architecture
func GetOsArch() string {
	return OsArch
}
