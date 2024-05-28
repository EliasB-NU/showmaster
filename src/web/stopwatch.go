package web

import (
	"sync"
	"time"
)

type Stopwatch struct {
	startTime   time.Time
	accumulated time.Duration
	running     bool
	mutex       sync.Mutex
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{}
}

func (s *Stopwatch) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.running {
		s.startTime = time.Now()
		s.running = true
		go s.runTimer()
	}
}

func (s *Stopwatch) Stop() time.Duration {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		s.accumulated += time.Since(s.startTime)
		s.running = false
	}
	return s.accumulated
}

func (s *Stopwatch) Resume() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.running {
		s.startTime = time.Now()
		s.running = true
		go s.runTimer()
	}
}

func (s *Stopwatch) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.running = false
	s.accumulated = 0
}

func (s *Stopwatch) ElapsedSeconds() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		return s.accumulated.Seconds() + time.Since(s.startTime).Seconds()
	}
	return s.accumulated.Seconds()
}

func (s *Stopwatch) runTimer() {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		if !s.running {
			s.mutex.Unlock()
			return
		}
		s.mutex.Unlock()
	}
}
