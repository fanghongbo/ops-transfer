package main

import "github.com/fanghongbo/ops-transfer/common/g"

var (
	Version    = "v1.0"
	BinaryName = "ops-transfer"
)

func init() {
	g.BinaryName = BinaryName
	g.Version = Version
}
