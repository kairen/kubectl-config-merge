package version

// version will be overridden with the current version at build time using the -X linker flag
var version = "v0.0.0-unset"

func GetVersion() string {
	return version
}
