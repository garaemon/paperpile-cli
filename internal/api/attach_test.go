package api

import (
	"testing"
)

func TestComputeMD5(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "empty input",
			input: []byte{},
			want:  "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:  "hello world",
			input: []byte("hello world"),
			want:  "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
		{
			name:  "binary data",
			input: []byte{0x00, 0x01, 0x02, 0xff},
			want:  "0416dab819887333af831f8c765ac2ae",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeMD5(tt.input)
			if got != tt.want {
				t.Errorf("computeMD5() = %q, want %q", got, tt.want)
			}
		})
	}
}
