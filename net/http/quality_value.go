package httpext

import "fmt"

const (
	// QualityValueFormat is a format string helper for Quality Values
	QualityValueFormat = "%s;q=%1.3g"
)

// QualityValue accepts a values to add/concatenate a quality values to and
// the quality values itself.
func QualityValue(v string, qv float32) string {
	if qv > 1 {
		qv = 1 // highest possible values
	}
	if qv < 0.001 {
		qv = 0.001 // lowest possible values
	}
	return fmt.Sprintf(QualityValueFormat, v, qv)
}
