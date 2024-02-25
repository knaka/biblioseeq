package common

import "testing"

func TestFoo(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"First", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Foo(); (err != nil) != tt.wantErr {
				t.Errorf("Foo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
