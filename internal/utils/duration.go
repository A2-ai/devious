package utils

import "time"

// FormatDuration formats a duration with a precision of 3 digits
// if it is less than 100s.
func FormatDuration(d time.Duration) string {
	scale := 100 * time.Second
	// look for the max scale that is smaller than d
	for scale > d {
		scale = scale / 10
	}
	return d.Round(scale / 100).String()
}
