//go:build !windows
// +build !windows

package main

import "github.com/warpbuilds/warpbuild-agent/cmd/agentd/cmd"

func main() {
	cmd.Execute()
}
