package plugin

import (
	"fmt"
	"time"
)

const (
	SleepCmd1 = "SleepCmd1"

	name1 = "SleepCmd1"
)

var (
	duration1 = 3 * time.Second
)

type Sleep1 struct {
	Name     string
	Duration time.Duration
}

func NewSleep1() *Sleep1 {
	return &Sleep1{Name: name1, Duration: duration1}
}

func (s *Sleep1) Run() error {
	fmt.Println("Name", s.Name, "sleeping...")
	time.Sleep(s.Duration)
	return nil
}

func (s *Sleep1) Stop() error {
	fmt.Println("Name", s.Name, "stopped")
	return nil
}
