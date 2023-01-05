package timeext

const (
	// RFC3339Nano is a correct replacement to Go's current time.RFC3339Nano which is NOT sortable and
	// have no intention of fixing that https://github.com/golang/go/issues/19635; this format fixes that.
	RFC3339Nano = "2006-01-02T15:04:05.000000000Z07:00"
)
