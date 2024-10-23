package chans

// Closed returns whether the chan is closed.
func Closed[T any](ch <-chan T) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

// CloseOnce close the chan if it is not closed.
// It does nothing when the chan is closed.
func CloseOnce[T any](ch chan T) {
	if !Closed(ch) {
		close(ch)
	}
}

// Done returns whether the chan is closed.
// It is always used to judge whether context.Context is done.
// eg. if Done(ctx.Done()) { ... }
func Done(done <-chan struct{}) bool {
	return Closed(done)
}
