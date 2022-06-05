package plugin

import (
	"testing"
	"time"
)

func TestSleep2_Run(t *testing.T) {
	type fields struct {
		Name     string
		Duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Sleep2 run",
			fields: fields{
				Name:     "Sleep 2",
				Duration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sleep2{
				Name:     tt.fields.Name,
				Duration: tt.fields.Duration,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Sleep2.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSleep2_Stop(t *testing.T) {
	type fields struct {
		Name     string
		Duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Sleep2 stop",
			fields: fields{
				Name:     "Sleep 2",
				Duration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sleep2{
				Name:     tt.fields.Name,
				Duration: tt.fields.Duration,
			}
			if err := s.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Sleep2.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
