package session

import "sync"

type Queue struct {
	mu    sync.Mutex
	items []uint64
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]uint64, 0, 1024),
	}
}

func (s *Queue) Push(v uint64) {
	s.mu.Lock()
	s.items = append(s.items, v)
	s.mu.Unlock()
}

// Pop returns the oldest pushed item.
func (s *Queue) Pop() (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := len(s.items)
	if n == 0 {
		return 0, false
	}

	v := s.items[0]
	copy(s.items, s.items[1:])
	s.items[n-1] = 0
	s.items = s.items[:n-1]

	return v, true
}

func (s *Queue) Len() int {
	s.mu.Lock()
	n := len(s.items)
	s.mu.Unlock()
	return n
}
