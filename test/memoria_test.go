package test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	memoria "github.com/IMGIITRoorkee/Memoria_Simple"
)

func TestMemoriaWriteRead(t *testing.T) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "memoria-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	m := memoria.New(memoria.Options{
		Basedir:      tempDir,
		MaxCacheSize: 1024,
	})

	tests := []struct {
		name    string
		key     string
		value   []byte
		wantErr bool
	}{
		{
			name:    "Simple write and read",
			key:     "test1",
			value:   []byte("hello world"),
			wantErr: false,
		},
		{
			name:    "Empty key",
			key:     "",
			value:   []byte("test"),
			wantErr: true,
		},
		{
			name:    "Empty value",
			key:     "test2",
			value:   []byte{},
			wantErr: false,
		},
		{
			name:    "Large value",
			key:     "test3",
			value:   bytes.Repeat([]byte("a"), 1000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Write
			err := m.Write(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Test Read
			got, err := m.Read(tt.key)
			if err != nil {
				t.Errorf("Read() error = %v", err)
				return
			}

			if !bytes.Equal(got, tt.value) {
				t.Errorf("Read() got = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestMemoriaWriteReadString(t *testing.T) {

	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "memoria-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	m := memoria.New(memoria.Options{
		Basedir:      tempDir,
		MaxCacheSize: 1024,
	})

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "Simple write and read (String Wrapper)",
			key:     "test1",
			value:   string("Hello World!"),
			wantErr: false,
		},
		{
			name:    "Empty key",
			key:     "",
			value:   string("test"),
			wantErr: true,
		},
		{
			name:    "Empty value",
			key:     "test2",
			value:   string(""),
			wantErr: false,
		},
		{
			name:    "Large value",
			key:     "test3",
			value:   strings.Repeat("a", 1000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test WriteString
			err := m.WriteString(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Test ReadString
			got, err := m.ReadString(tt.key)
			if err != nil {
				t.Errorf("Read() error = %v", err)
				return
			}

			if got != tt.value {
				t.Errorf("Read() got = %v, want %v", got, tt.value)
			}
		})
	}
}
