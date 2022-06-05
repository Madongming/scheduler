package plugin

import (
	"testing"
	"time"
)

func TestSleep1_Run(t *testing.T) {
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
			name: "Sleep1 run",
			fields: fields{
				Name:     "Sleep 1",
				Duration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sleep1{
				Name:     tt.fields.Name,
				Duration: tt.fields.Duration,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Sleep1.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSleep1_Stop(t *testing.T) {
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
			name: "Sleep1 stop",
			fields: fields{
				Name:     "Sleep 1",
				Duration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sleep1{
				Name:     tt.fields.Name,
				Duration: tt.fields.Duration,
			}
			if err := s.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Sleep1.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
