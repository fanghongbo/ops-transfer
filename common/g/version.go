package g

import "fmt"

var (
	Version    string = "v0.1"
	BinaryName string = "ops-transfer.dev"
)

func VersionInfo() string {
	return fmt.Sprintf("%s", Version)
}

func HbsInfo() string {
	return fmt.Sprintf("%s.%s", BinaryName, Version)
}
