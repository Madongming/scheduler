package plugin

import (
	"fmt"
	"time"
)

const (
	SleepCmd2 = "SleepCmd2"

	name2 = "SleepCmd2"
)

var (
	duration2 = 3 * time.Second
)

type Sleep2 struct {
	Name     string
	Duration time.Duration
}

func NewSleep2() *Sleep2 {
	return &Sleep2{Name: name2, Duration: duration2}
}

func (s *Sleep2) Run() error {
	fmt.Println("Name", s.Name, "sleeping...")
	time.Sleep(s.Duration)
	return nil
}

func (s *Sleep2) Stop() error {
	fmt.Println("Name", s.Name, "stopped")
	return nil
}
