package errorsext

import (
	"errors"
	"strings"
	"syscall"
)

var (
	// ErrMaxAttemptsReached is a placeholder error to use when some retryable even has reached its maximum number of
	// attempts.
	ErrMaxAttemptsReached = errors.New("max attempts reached")
)

// IsRetryableHTTP returns if the provided error is considered retryable HTTP error. It also returns the
// type, in string form, for optional logging and metrics use.
func IsRetryableHTTP(err error) (retryType string, isRetryable bool) {
	if retryType, isRetryable = IsRetryableNetwork(err); isRetryable {
		return
	}

	errStr := err.Error()

	if strings.Contains(errStr, "http2: server sent GOAWAY") {
		return "goaway", true
	}

	// errServerClosedIdle is not seen by users for idempotent HTTP requests, but may be
	// seen by a user if the server shuts down an idle connection and sends its FIN
	// in flight with already-written POST body bytes from the client.
	// See https://github.com/golang/go/issues/19943#issuecomment-355607646
	//
	// This will possibly get fixed in the upstream SDK's based on the ability to set an HTTP error in the future
	// https://go-review.googlesource.com/c/go/+/191779/ but until then we should retry these.
	//
	if strings.Contains(errStr, "http: server closed idle connection") {
		return "server_close_idle_connection", true
	}
	return "", false
}

// IsRetryableNetwork returns if the provided error is a retryable network related error. It also returns the
// type, in string form, for optional logging and metrics use.
func IsRetryableNetwork(err error) (retryType string, isRetryable bool) {
	if IsRetryable(err) {
		return "retryable", true
	}
	if IsTemporary(err) {
		return "temporary", true
	}
	if IsTimeout(err) {
		return "timeout", true
	}
	return IsTemporaryConnection(err)
}

// IsRetryable returns true if the provided error is considered retryable by testing if it
// complies with an interface implementing `Retryable() bool` or `IsRetryable bool` and calling the function.
func IsRetryable(err error) bool {
	var t interface{ IsRetryable() bool }
	if errors.As(err, &t) && t.IsRetryable() {
		return true
	}

	var t2 interface{ Retryable() bool }
	return errors.As(err, &t2) && t2.Retryable()
}

// IsTemporary returns true if the provided error is considered retryable temporary error by testing if it
// complies with an interface implementing `Temporary() bool` and calling the function.
func IsTemporary(err error) bool {
	var t interface{ Temporary() bool }
	return errors.As(err, &t) && t.Temporary()
}

// IsTimeout returns true if the provided error is considered a retryable timeout error by testing if it
// complies with an interface implementing `Timeout() bool` and calling the function.
func IsTimeout(err error) bool {
	var t interface{ Timeout() bool }
	return errors.As(err, &t) && t.Timeout()
}

// IsTemporaryConnection returns if the provided error was a low level retryable connection error. It also returns the
// type, in string form, for optional logging and metrics use.
func IsTemporaryConnection(err error) (retryType string, isRetryable bool) {
	if err != nil {
		if errors.Is(err, syscall.ECONNRESET) {
			return "econnreset", true
		}
		if errors.Is(err, syscall.ECONNABORTED) {
			return "econnaborted", true
		}
		if errors.Is(err, syscall.ENOTCONN) {
			return "enotconn", true
		}
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return "ewouldblock", true
		}
		if errors.Is(err, syscall.EAGAIN) {
			return "eagain", true
		}
		if errors.Is(err, syscall.ETIMEDOUT) {
			return "etimedout", true
		}
		if errors.Is(err, syscall.EINTR) {
			return "eintr", true
		}
		if errors.Is(err, syscall.EPIPE) {
			return "epipe", true
		}
	}
	return "", false
}
