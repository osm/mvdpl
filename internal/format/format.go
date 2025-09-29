package format

import "fmt"

func Time(seconds float64) string {
	totalSeconds := int(seconds)
	minutes := totalSeconds / 60
	secs := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}
