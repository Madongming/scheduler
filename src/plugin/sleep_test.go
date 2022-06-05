package plugin

import (
	"testing"
	"time"
)

func TestSleep_Run(t *testing.T) {
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
			name: "Sleep run",
			fields: fields{
				Name:     "Sleep 1",
				Duration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sleep{
				Name:     tt.fields.Name,
				Duration: tt.fields.Duration,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Sleep.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSleep_Stop(t *testing.T) {
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
			name: "Sleep stop",
			fields: fields{
				Name:     "Sleep 1",
				Duration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sleep{
				Name:     tt.fields.Name,
				Duration: tt.fields.Duration,
			}
			if err := s.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Sleep.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
