package main

import (
	"dvs/cmd"
	"dvs/internal"
)

func main() {
	// Clear default umask
	internal.ClearUmask()

	cmd.Execute()
}
