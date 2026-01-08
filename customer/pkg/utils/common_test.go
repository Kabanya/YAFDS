package utils

import (
	"runtime"
	"testing"
)

func TestNumThreads(t *testing.T) {
	maxProcs := runtime.GOMAXPROCS(0)
	tests := []struct {
		name      string
		requested int
		want      uint8
	}{
		{
			name:      "requested zero",
			requested: 0,
			want:      uint8(maxProcs),
		},
		{
			name:      "requested negative",
			requested: -1,
			want:      uint8(maxProcs),
		},
		{
			name:      "requested one",
			requested: 1,
			want:      1,
		},
		{
			name:      "requested more than max",
			requested: maxProcs + 1,
			want:      uint8(maxProcs),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumThreads(tt.requested); got != tt.want {
				t.Errorf("NumThreads(%d) = %v, want %v", tt.requested, got, tt.want)
			}
		})
	}
}
