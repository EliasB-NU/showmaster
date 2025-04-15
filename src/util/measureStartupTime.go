package util

import (
	"fmt"
	"time"
)

type MST struct {
	time time.Time
}

func (m *MST) StartTimer() {
	m.time = time.Now()
}

func (m *MST) ElapsedTime() {
	fmt.Printf("Total startup time : %v\n", time.Since(m.time))
}
