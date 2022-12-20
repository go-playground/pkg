package contextext

import (
	"context"
	"fmt"
	"time"
)

var _ context.Context = (*detachedContext)(nil)

type detachedContext struct {
	parent context.Context
}

// Detach returns a new context which continues to have access to its parent values but
// is no longer bound/attached to the parents timeouts nor deadlines.
//
// If a nil context is passed in a new Background context will be used as the parent.
//
// This is useful for when you wish to pass along values such as span or logging information but do not want the
// current operation to be cancelled despite what upstream code/callers think.
func Detach(parent context.Context) detachedContext {
	if parent == nil {
		return detachedContext{parent: context.Background()}
	}
	return detachedContext{parent: parent}
}

func (c detachedContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c detachedContext) Done() <-chan struct{} {
	return nil
}

func (c detachedContext) Err() error {
	return nil
}

func (c detachedContext) Value(key any) any {
	return c.parent.Value(key)
}

func (c detachedContext) String() string {
	return fmt.Sprintf("%s.Detached", c.parent)
}
