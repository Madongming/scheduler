package plugin

import (
	"fmt"
	"time"
)

type Sleep struct {
	Name     string
	Duration time.Duration
}

func NewSleep(name string, duration time.Duration) *Sleep {
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
