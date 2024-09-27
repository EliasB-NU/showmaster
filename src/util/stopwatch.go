package util

import (
	"sync"
	"time"
)

type Stopwatch struct {
	startTime   time.Time
	accumulated time.Duration
	Running     bool
	mutex       sync.Mutex
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{}
}

func (s *Stopwatch) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.Running {
		s.startTime = time.Now()
		s.Running = true
		go s.runTimer()
	}
}

func (s *Stopwatch) Stop() time.Duration {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.Running {
		s.accumulated += time.Since(s.startTime)
		s.Running = false
	}
	return s.accumulated
}

func (s *Stopwatch) Resume() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.Running {
		s.startTime = time.Now()
		s.Running = true
		go s.runTimer()
	}
}

func (s *Stopwatch) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Running = false
	s.accumulated = 0
}

func (s *Stopwatch) ElapsedSeconds() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.Running {
		return s.accumulated.Seconds() + time.Since(s.startTime).Seconds()
	}
	return s.accumulated.Seconds()
}

func (s *Stopwatch) runTimer() {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		if !s.Running {
			s.mutex.Unlock()
			return
		}
		s.mutex.Unlock()
	}
}

func (s *Stopwatch) SetElapsedSeconds(elapsed float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.accumulated += time.Duration(elapsed)
}
