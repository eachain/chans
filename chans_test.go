package chans

import (
	"context"
	"testing"
)

func TestClosed(t *testing.T) {
	ch := make(chan struct{})
	if Closed(ch) {
		t.Fatalf("chan not closed, but returns closed")
	}
	close(ch)
	if !Closed(ch) {
		t.Fatalf("chan closed, but returns not closed")
	}
}

func TestCloseOnce(t *testing.T) {
	ch := make(chan struct{})
	CloseOnce(ch)
	CloseOnce(ch) // should not panic
}

func TestDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	if Done(ctx.Done()) {
		t.Fatalf("context not done, but returns done")
	}
	cancel()
	if !Done(ctx.Done()) {
		t.Fatalf("context done, but returns not done")
	}
}
