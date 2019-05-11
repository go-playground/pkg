package httpext

import "fmt"

const (
	// QualityValueFormat is a format string helper for Quality Values
	QualityValueFormat = "%s;q=%1.3g"
)

// QualityValue accepts a value to add/concatenate a quality value to and
// the quality value itself.
func QualityValue(v string, qv float32) string {
	if qv > 1 {
		qv = 1 // highest possible value
	}
	if qv < 0.001 {
		qv = 0.001 // lowest possible value
	}
	return fmt.Sprintf(QualityValueFormat, v, qv)
}
