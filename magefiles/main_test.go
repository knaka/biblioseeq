package main

import (
	"os"
	"testing"
)

func TestPb_Gen(t *testing.T) {
	_ = os.Chdir("..")
	tests := []struct {
		name    string
		pb      PB
		wantErr bool
	}{
		{
			"Test 1",
			PB{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pb.Gen(); (err != nil) != tt.wantErr {
				t.Errorf("Gen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
