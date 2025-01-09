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

func TestMemoriaWriteReadString(t1 *testing.T) {

	// Create temporary directory for tests
	tempDir1, err1 := os.MkdirTemp("", "memoria-test-*")
	if err1 != nil {
		t1.Fatalf("Failed to create temp dir: %v", err1)
	}
	defer os.RemoveAll(tempDir1)

	m := memoria.New(memoria.Options{
		Basedir:      tempDir1,
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

	for _, tt1 := range tests {
		t1.Run(tt1.name, func(t1 *testing.T) {

			// Test WriteString
			err1 := m.WriteString(tt1.key, tt1.value)
			if (err1 != nil) != tt1.wantErr {
				t1.Errorf("Write() error = %v, wantErr %v", err1, tt1.wantErr)
				return
			}

			if tt1.wantErr {
				return
			}

			// Test ReadString
			got1, err1 := m.ReadString(tt1.key)
			if err1 != nil {
				t1.Errorf("Read() error = %v", err1)
				return
			}

			if got1 != tt1.value {
				t1.Errorf("Read() got = %v, want %v", got1, tt1.value)
			}
		})
	}
}

func TestMemoriaBulkWrite(t *testing.T) {

	// temp-dir for tests
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
		name        string
		pairs       map[string][]byte
		numWorkers  int
		expectedErr bool
	}{
		{
			name: "Successful bulk write with multiple pairs",
			pairs: map[string][]byte{
				"key1": []byte("value1"),
				"key2": []byte("value2"),
				"key3": []byte("value3"),
			},
			numWorkers:  2,
			expectedErr: false,
		},
		{
			name: "Bulk write with an empty key",
			pairs: map[string][]byte{
				"": []byte("value1"), // Empty key should cause error
			},
			numWorkers:  1,
			expectedErr: true,
		},
		{
			name: "Bulk write with an empty value",
			pairs: map[string][]byte{
				"key1": []byte{}, // Empty value should be handled correctly
			},
			numWorkers:  1,
			expectedErr: false,
		},
		{
			name: "Bulk write with large values",
			pairs: map[string][]byte{
				"key1": bytes.Repeat([]byte("a"), 1000), // Large value
				"key2": bytes.Repeat([]byte("b"), 1000), // Large value
			},
			numWorkers:  2,
			expectedErr: false,
		},
		{
			name: "Bulk write with mix of valid and invalid keys",
			pairs: map[string][]byte{
				"validKey1": []byte("value1"),
				"":          []byte("value2"), // invalid key
			},
			numWorkers:  2,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			results := m.BulkWrite(tt.pairs, tt.numWorkers)

			// Check the number of results matches the number of pairs
			if len(results) != len(tt.pairs) {
				t.Errorf("Expected %d results, but got %d", len(tt.pairs), len(results))
			}

			var hasError bool
			for _, result := range results {
				if result.Error != nil {
					hasError = true
				}
			}

			// Check if the results contain errors where expected
			// for _, result := range results {
			// 	if tt.expectedErr && result.Error == nil {
			// 		t.Errorf("Expected error but got nil for key %s", result.Key)
			// 	}
			// 	if !tt.expectedErr && result.Error != nil {
			// 		t.Errorf("Unexpected error for key %s: %v", result.Key, result.Error)
			// 	}
			// }

			// If we expect an error, ensure at least one result has an error
			if tt.expectedErr && !hasError {
				t.Errorf("Expected at least one error but got none")
			}

			// If no error is expected, ensure all results are successful
			if !tt.expectedErr && hasError {
				t.Errorf("Unexpected error(s) occurred")
			}

			// Verify the data was written correctly
			for key, expectedValue := range tt.pairs {
				if tt.expectedErr {
					continue
				}
				got, err := m.Read(key)
				if err != nil {
					t.Errorf("Failed to read key %s: %v", key, err)
					continue
				}
				if !bytes.Equal(got, expectedValue) {
					t.Errorf("Read value for key %s = %v, want %v", key, got, expectedValue)
				}
			}
		})
	}
}
