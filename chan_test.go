package chans

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	ch := New[int]()
	ch.Send(123)
	ch.Send(456)
	ch.Send(789)
}

func TestSendNilChan(t *testing.T) {
	done := make(chan struct{})
	go func() {
		var ch *Chan[int]
		ch.Send(123) // it will blocked
		close(done)
	}()

	select {
	case <-time.After(10 * time.Millisecond):
	case <-done:
		t.Fatalf("send to nil chan done")
	}
}

func TestSendClosedChan(t *testing.T) {
	ch := New[int]()
	ch.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("send to closed chan not panic")
		}
	}()
	ch.Send(123)
}

func TestRecv(t *testing.T) {
	ch := New[int]()

	x := 123
	ch.Send(x)
	if v := ch.Recv(); v != x {
		t.Fatalf("chan recv data: %v, send: %v", v, x)
	}

	// async send
	y := 456
	z := 789
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch.Send(y)
		time.Sleep(10 * time.Millisecond)
		ch.Send(z)
	}()

	if v := ch.Recv(); v != y {
		t.Fatalf("chan recv data: %v, send: %v", v, x)
	}
	if v := ch.Recv(); v != z {
		t.Fatalf("chan recv data: %v, send: %v", v, x)
	}
}

func TestRecvBlockClose(t *testing.T) {
	ch := New[int]()

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch.Close()
	}()

	if v := ch.Recv(); v != 0 {
		t.Fatalf("chan recv data from closed chan: %v", v)
	}
}

func TestRecvNilChan(t *testing.T) {
	done := make(chan struct{})
	go func() {
		var ch *Chan[int]
		ch.Recv() // it will blocked
		close(done)
	}()

	select {
	case <-time.After(10 * time.Millisecond):
	case <-done:
		t.Fatalf("recv from nil chan done")
	}
}

func TestRecvClosed(t *testing.T) {
	ch := New[int]()
	x := 123
	ch.Send(x)
	ch.Close()

	v, ok := ch.TryRecv()
	if !ok {
		t.Fatalf("chan recv not ok")
	}
	if v != x {
		t.Fatalf("chan recv data: %v, send: %v", v, x)
	}

	v, ok = ch.TryRecv()
	if ok {
		t.Fatalf("chan recv closed chan ok")
	}
	if v != 0 {
		t.Fatalf("chan recv closed chan value: %v", v)
	}
}

func TestClose(t *testing.T) {
	ch := New[int]()
	wg := new(sync.WaitGroup)
	const N = 10
	errs := make([]string, N)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				v, ok := ch.TryRecv()
				if ok {
					errs[i] = fmt.Sprintf("go[%v] received: %v", i, v)
					continue
				}
				break
			}
		}(i)
	}

	ch.Close()
	wg.Wait()

	for _, err := range errs {
		if err != "" {
			t.Errorf(err)
		}
	}
}

func TestCloseNilChan(t *testing.T) {
	var ch *Chan[int]

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("close nil chan not panic")
		}
	}()

	ch.Close()
}

func TestCloseClosedChan(t *testing.T) {
	ch := New[int]()
	ch.Send(123)
	ch.Recv()
	ch.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("close closed chan not panic")
		}
	}()
	ch.Close() // double close will not panic
}
