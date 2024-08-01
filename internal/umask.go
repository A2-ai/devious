//go:build !windows

package internal

import "syscall"

func ClearUmask() {
	syscall.Umask(0)
}
