package storage

import "testing"

func TestLoadRegistry(t *testing.T) {
	err := LoadRegistary()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
