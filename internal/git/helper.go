package git

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// Removes lines from a file
// Adapted from https://rosettacode.org/wiki/Remove_lines_from_a_file#Go
func removeLines(file *os.File, start, n int) (err error) {
	if start < 1 {
		return errors.New("Invalid request. Line numbers start at 1.")
	}
	if n < 0 {
		return errors.New("invalid request. Negative number to remove.")
	}
	var b []byte
	if b, err = io.ReadAll(file); err != nil {
		return
	}
	cut, ok := skip(b, start-1)
	if !ok {
		return fmt.Errorf("less than %d lines", start)
	}
	if n == 0 {
		return nil
	}
	tail, ok := skip(cut, n)
	if !ok {
		return fmt.Errorf("less than %d lines after line %d", n, start)
	}
	t := int64(len(b) - len(cut))
	if err = file.Truncate(t); err != nil {
		return
	}
	if len(tail) > 0 {
		_, err = file.WriteAt(tail, t)
	}
	return
}

func skip(b []byte, n int) ([]byte, bool) {
	for ; n > 0; n-- {
		if len(b) == 0 {
			return nil, false
		}
		x := bytes.IndexByte(b, '\n')
		if x < 0 {
			x = len(b)
		} else {
			x++
		}
		b = b[x:]
	}
	return b, true
}
