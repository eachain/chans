package chans

import (
	"sync"
)

const maxFreeCache = 1000

type node[T any] struct {
	data T
	next *node[T]
}

// Chan[T] like builtin chan type, but without capacity limit.
type Chan[T any] struct {
	cond   sync.Cond
	head   *node[T]
	tail   *node[T]
	free   *node[T]
	freed  int
	closed bool
}

// New[T] like `make(chan T)` but without capacity parameter.
// It returns a *Chan[T] without capacity limit.
func New[T any]() *Chan[T] {
	ch := new(Chan[T])
	ch.cond.L = new(sync.Mutex)
	return ch
}

func (ch *Chan[T]) isNil() bool {
	return ch == nil || ch.cond.L == nil
}

// Send like `ch <- v`.
// It send a value to the *Chan.
// It will block when *Chan is nil.
func (ch *Chan[T]) Send(v T) {
	if ch.isNil() {
		select {} // chan send (nil chan)
	}
	ch.send(v)
	ch.cond.Signal()
}

func (ch *Chan[T]) send(v T) {
	ch.cond.L.Lock()
	defer ch.cond.L.Unlock()

	if ch.closed {
		panic("send on closed channel")
	}

	n := ch.free
	if n != nil {
		ch.free = n.next
		n.next = nil
		ch.freed--
	} else {
		n = new(node[T])
	}

	n.data = v

	if ch.tail != nil {
		ch.tail.next = n
		ch.tail = n
	} else {
		ch.head = n
		ch.tail = n
	}
}

// Recv like `v := <-ch`.
// It try to receive a value from *Chan.
// It returns a value when here is any cached value in *Chan.
// It returns a zero value
// when here is not any cached value in *Chan and the *Chan is closed.
// It will block when *Chan is nil.
func (ch *Chan[T]) Recv() (v T) {
	v, _ = ch.TryRecv()
	return
}

// TryRecv like `v, ok := <-ch`.
// It try to receive a value from *Chan.
// It returns a value and true when here is any cached value in *Chan.
// It returns a zero value and false
// when here is not any cached value in *Chan and the *Chan is closed.
// It will block when *Chan is nil.
func (ch *Chan[T]) TryRecv() (v T, ok bool) {
	if ch.isNil() {
		select {}
	}

	ch.cond.L.Lock()
	defer ch.cond.L.Unlock()

	for !ch.closed && ch.head == nil {
		ch.cond.Wait()
	}
	if ch.closed && ch.head == nil {
		return
	}

	ok = true
	v = ch.head.data

	var zero T
	ch.head.data = zero

	n := ch.head
	ch.head = ch.head.next
	if ch.head == nil {
		ch.tail = nil
	}

	if !ch.closed && ch.freed < maxFreeCache {
		n.next = ch.free
		ch.free = n
		ch.freed++
	} else {
		n.next = nil
	}
	return
}

// Close like builtin close(chan).
// It will panic if the *Chan is nil or the *Chan is closed.
func (ch *Chan[T]) Close() {
	if ch.isNil() {
		panic("close of nil channel")
	}

	ch.cond.L.Lock()

	if ch.closed {
		ch.cond.L.Unlock()
		panic("close of closed channel")
	}
	ch.closed = true

	var zero T
	var next *node[T]
	for n := ch.free; n != nil; n = next {
		next = n.next
		n.data = zero
		n.next = nil
	}
	ch.free = nil
	ch.freed = 0

	ch.cond.L.Unlock()

	ch.cond.Broadcast()
}
