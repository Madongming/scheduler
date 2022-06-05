package plugin

import (
	"fmt"
	"time"
)

const (
	SleepCmd = "SleepCmd"

	name = "SleepCmd"
)

var (
	duration = 3 * time.Second
)

type Sleep struct {
	Name     string
	Duration time.Duration
}

func NewSleep() *Sleep {
	return &Sleep{Name: name, Duration: duration}
}

func (s *Sleep) Run() error {
	fmt.Println("Name", s.Name, "sleeping...")
	time.Sleep(s.Duration)
	return nil
}

func (s *Sleep) Stop() error {
	fmt.Println("Name", s.Name, "stopped")
	return nil
}
