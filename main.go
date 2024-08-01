package main

import (
	"dvs/cmd"
	"syscall"
)

func main() {
	// Clear default umask
	syscall.Umask(0)

	cmd.Execute()
}
