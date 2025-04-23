package util

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

func (s *Stopwatch) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		s.accumulated += time.Since(s.startTime)
		s.running = false
	}
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

func (s *Stopwatch) ElapsedSeconds() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		return int64(s.accumulated.Seconds() + time.Since(s.startTime).Seconds())
	}
	return int64(s.accumulated.Seconds())
}

func (s *Stopwatch) setElapsedSeconds(elapsedSeconds int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		s.running = false
		s.accumulated = time.Duration(elapsedSeconds) * time.Second
	}
}

func (s *Stopwatch) runTimer() {
	ticker := time.NewTicker(time.Second)
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
