package contextext

import (
	"context"
	"testing"
	"time"
)

func TestDetach(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", 13)
	ctx, cancel := context.WithTimeout(ctx, time.Nanosecond)
	cancel() // cancel ensuring context has been canceled

	select {
	case <-ctx.Done():
	default:
		t.Fatal("expected context to be cancelled")
	}

	ctx = Detach(ctx)

	select {
	case <-ctx.Done():
		t.Fatal("expected context to be detached from parents cancellation")
	default:
		rawValue := ctx.Value("key")
		n, ok := rawValue.(int)
		if !ok || n != 13 {
			t.Fatalf("expected integer woth value of 13 but got %v", rawValue)
		}
	}

	// test new detached nil parent context
	ctx = Detach(nil)
	select {
	case <-ctx.Done():
		t.Fatal("expected context to be detached from parents cancellation")
	default:
	}
}
