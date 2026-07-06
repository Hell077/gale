package session

import "testing"

func TestQueuePopIsFIFO(t *testing.T) {
	q := NewQueue()
	q.Push(1)
	q.Push(2)
	q.Push(3)

	for _, want := range []uint64{1, 2, 3} {
		got, ok := q.Pop()
		if !ok {
			t.Fatalf("expected %d, queue was empty", want)
		}
		if got != want {
			t.Fatalf("expected %d, got %d", want, got)
		}
	}

	if got, ok := q.Pop(); ok {
		t.Fatalf("expected empty queue, got %d", got)
	}
}
