package config

import (
	"testing"
)

func TestCreateConfigFile(t *testing.T) {
	var want error
	wantErr := false
	err := CreateConfigFile()
	if err != nil && !wantErr {
		t.Errorf("Helper() = %v, want %v", err, want)
	}
}
