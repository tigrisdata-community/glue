package store

import (
	"context"
	"errors"
	"testing"
)

// mockStore is an in-memory implementation for testing.
type mockStore struct {
	data map[string][]byte
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[string][]byte),
	}
}

func (m *mockStore) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *mockStore) Exists(ctx context.Context, key string) error {
	if _, ok := m.data[key]; !ok {
		return ErrNotFound
	}
	return nil
}

func (m *mockStore) Get(ctx context.Context, key string) ([]byte, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, ErrNotFound
	}
	return val, nil
}

func (m *mockStore) Set(ctx context.Context, key string, value []byte) error {
	m.data[key] = value
	return nil
}

func (m *mockStore) List(ctx context.Context, prefix string) ([]string, error) {
	var result []string
	for k := range m.data {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			result = append(result, k)
		}
	}
	return result, nil
}

func TestInterface_Exists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockStore)
		key      string
		wantErr  error
		errCheck func(error) bool
	}{
		{
			name: "returns nil when key exists",
			setup: func(m *mockStore) {
				m.data["existing-key"] = []byte("value")
			},
			key:     "existing-key",
			wantErr: nil,
		},
		{
			name:    "returns ErrNotFound when key does not exist",
			setup:   func(m *mockStore) {},
			key:     "non-existent-key",
			wantErr: ErrNotFound,
			errCheck: func(err error) bool {
				return errors.Is(err, ErrNotFound)
			},
		},
		{
			name: "returns nil for empty key that exists",
			setup: func(m *mockStore) {
				m.data[""] = []byte("empty-key-value")
			},
			key:     "",
			wantErr: nil,
		},
		{
			name: "returns ErrNotFound for empty key that does not exist",
			setup: func(m *mockStore) {
				m.data["other-key"] = []byte("value")
			},
			key:     "",
			wantErr: ErrNotFound,
			errCheck: func(err error) bool {
				return errors.Is(err, ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			err := m.Exists(context.Background(), tt.key)

			if tt.wantErr != nil {
				if tt.errCheck != nil {
					if !tt.errCheck(err) {
						t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
					}
				} else if err != tt.wantErr {
					t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("Exists() unexpected error = %v", err)
			}
		})
	}
}

func TestJSON_Exists(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		setup    func(*mockStore)
		key      string
		wantErr  error
		errCheck func(error) bool
	}{
		{
			name: "returns nil when key exists with prefix",
			setup: func(m *mockStore) {
				m.data["testprefix/mykey"] = []byte(`{"value":"data"}`)
			},
			prefix:  "testprefix",
			key:     "mykey",
			wantErr: nil,
		},
		{
			name:    "returns ErrNotFound when key does not exist with prefix",
			setup:   func(m *mockStore) {},
			prefix:  "testprefix",
			key:     "mykey",
			wantErr: ErrNotFound,
			errCheck: func(err error) bool {
				return errors.Is(err, ErrNotFound)
			},
		},
		{
			name: "returns nil when key exists without prefix",
			setup: func(m *mockStore) {
				m.data["mykey"] = []byte(`{"value":"data"}`)
			},
			prefix:  "",
			key:     "mykey",
			wantErr: nil,
		},
		{
			name: "returns nil when key exists with empty prefix",
			setup: func(m *mockStore) {
				m.data["mykey"] = []byte(`{"value":"data"}`)
			},
			prefix:  "",
			key:     "mykey",
			wantErr: nil,
		},
		{
			name: "returns ErrNotFound when prefixed key not in store",
			setup: func(m *mockStore) {
				m.data["otherprefix/mykey"] = []byte(`{"value":"data"}`)
			},
			prefix:  "testprefix",
			key:     "mykey",
			wantErr: ErrNotFound,
			errCheck: func(err error) bool {
				return errors.Is(err, ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			j := &JSON[struct{ Value string }]{
				Underlying: m,
				Prefix:     tt.prefix,
			}

			err := j.Exists(context.Background(), tt.key)

			if tt.wantErr != nil {
				if tt.errCheck != nil {
					if !tt.errCheck(err) {
						t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
					}
				} else if err != tt.wantErr {
					t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("Exists() unexpected error = %v", err)
			}
		})
	}
}

func TestInterface_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockStore)
		key     string
		wantErr error
	}{
		{
			name: "deletes existing key",
			setup: func(m *mockStore) {
				m.data["delete-me"] = []byte("value")
			},
			key:     "delete-me",
			wantErr: nil,
		},
		{
			name:    "deleting non-existent key is no-op",
			setup:   func(m *mockStore) {},
			key:     "non-existent",
			wantErr: nil,
		},
		{
			name: "deletes empty key",
			setup: func(m *mockStore) {
				m.data[""] = []byte("empty-key")
			},
			key:     "",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			err := m.Delete(context.Background(), tt.key)

			if err != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify key is gone
			if _, exists := m.data[tt.key]; exists {
				t.Errorf("Delete() key %q still exists in store", tt.key)
			}
		})
	}
}

func TestInterface_Get(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockStore)
		key      string
		want     []byte
		wantErr  error
		errCheck func(error) bool
	}{
		{
			name: "returns value for existing key",
			setup: func(m *mockStore) {
				m.data["mykey"] = []byte("myvalue")
			},
			key:     "mykey",
			want:    []byte("myvalue"),
			wantErr: nil,
		},
		{
			name:    "returns ErrNotFound for non-existent key",
			setup:   func(m *mockStore) {},
			key:     "non-existent",
			want:    nil,
			wantErr: ErrNotFound,
			errCheck: func(err error) bool {
				return errors.Is(err, ErrNotFound)
			},
		},
		{
			name: "returns empty byte slice for key with empty value",
			setup: func(m *mockStore) {
				m.data["empty-val"] = []byte{}
			},
			key:     "empty-val",
			want:    []byte{},
			wantErr: nil,
		},
		{
			name: "returns binary data correctly",
			setup: func(m *mockStore) {
				m.data["binary"] = []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0xfd}
			},
			key:     "binary",
			want:    []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0xfd},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			got, err := m.Get(context.Background(), tt.key)

			if tt.wantErr != nil {
				if tt.errCheck != nil {
					if !tt.errCheck(err) {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					}
				} else if err != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Get() unexpected error = %v", err)
				}
			}

			if tt.want != nil && !equalBytes(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterface_Set(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   []byte
		setup   func(*mockStore)
		wantErr error
		verify  func(*testing.T, *mockStore)
	}{
		{
			name:    "sets new key",
			key:     "new-key",
			value:   []byte("new-value"),
			setup:   func(m *mockStore) {},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got, ok := m.data["new-key"]
				if !ok {
					t.Error("Set() key was not stored")
				}
				if !equalBytes(got, []byte("new-value")) {
					t.Errorf("Set() stored value = %v, want %v", got, []byte("new-value"))
				}
			},
		},
		{
			name:  "overwrites existing key",
			key:   "existing-key",
			value: []byte("new-value"),
			setup: func(m *mockStore) {
				m.data["existing-key"] = []byte("old-value")
			},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got := m.data["existing-key"]
				if !equalBytes(got, []byte("new-value")) {
					t.Errorf("Set() stored value = %v, want %v", got, []byte("new-value"))
				}
			},
		},
		{
			name:    "sets empty key",
			key:     "",
			value:   []byte("empty-key-value"),
			setup:   func(m *mockStore) {},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got, ok := m.data[""]
				if !ok {
					t.Error("Set() empty key was not stored")
				}
				if !equalBytes(got, []byte("empty-key-value")) {
					t.Errorf("Set() stored value = %v, want %v", got, []byte("empty-key-value"))
				}
			},
		},
		{
			name:    "sets empty value",
			key:     "empty-value-key",
			value:   []byte{},
			setup:   func(m *mockStore) {},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got, ok := m.data["empty-value-key"]
				if !ok {
					t.Error("Set() key was not stored")
				}
				if len(got) != 0 {
					t.Errorf("Set() stored value length = %d, want 0", len(got))
				}
			},
		},
		{
			name:    "sets binary data",
			key:     "binary-key",
			value:   []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0xfd},
			setup:   func(m *mockStore) {},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got := m.data["binary-key"]
				want := []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0xfd}
				if !equalBytes(got, want) {
					t.Errorf("Set() stored value = %v, want %v", got, want)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			err := m.Set(context.Background(), tt.key, tt.value)

			if err != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.verify != nil {
				tt.verify(t, m)
			}
		})
	}
}

func TestInterface_List(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockStore)
		prefix  string
		want    []string
		wantErr error
	}{
		{
			name: "lists keys with matching prefix",
			setup: func(m *mockStore) {
				m.data["foo/a"] = []byte("1")
				m.data["foo/b"] = []byte("2")
				m.data["foo/c"] = []byte("3")
				m.data["bar/a"] = []byte("4")
			},
			prefix:  "foo/",
			want:    []string{"foo/a", "foo/b", "foo/c"},
			wantErr: nil,
		},
		{
			name: "returns empty list for non-existent prefix",
			setup: func(m *mockStore) {
				m.data["other/key"] = []byte("value")
			},
			prefix:  "noprefix/",
			want:    []string{},
			wantErr: nil,
		},
		{
			name: "lists all keys with empty prefix",
			setup: func(m *mockStore) {
				m.data["a"] = []byte("1")
				m.data["b"] = []byte("2")
				m.data["c"] = []byte("3")
			},
			prefix:  "",
			want:    []string{"a", "b", "c"},
			wantErr: nil,
		},
		{
			name:    "returns empty list for empty store",
			setup:   func(m *mockStore) {},
			prefix:  "",
			want:    []string{},
			wantErr: nil,
		},
		{
			name: "lists keys with single character prefix",
			setup: func(m *mockStore) {
				m.data["a1"] = []byte("1")
				m.data["a2"] = []byte("2")
				m.data["b1"] = []byte("3")
			},
			prefix:  "a",
			want:    []string{"a1", "a2"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			got, err := m.List(context.Background(), tt.prefix)

			if err != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !equalStringSlicesUnordered(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_Get(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name     string
		prefix   string
		setup    func(*mockStore)
		key      string
		want     testStruct
		wantErr  error
		errCheck func(error) bool
	}{
		{
			name: "unmarshals existing key with prefix",
			setup: func(m *mockStore) {
				m.data["testprefix/mykey"] = []byte(`{"name":"test","value":42}`)
			},
			prefix:  "testprefix",
			key:     "mykey",
			want:    testStruct{Name: "test", Value: 42},
			wantErr: nil,
		},
		{
			name:    "returns ErrNotFound for non-existent key with prefix",
			setup:   func(m *mockStore) {},
			prefix:  "testprefix",
			key:     "mykey",
			wantErr: ErrNotFound,
			errCheck: func(err error) bool {
				return errors.Is(err, ErrNotFound)
			},
		},
		{
			name: "unmarshals existing key without prefix",
			setup: func(m *mockStore) {
				m.data["mykey"] = []byte(`{"name":"noprefix","value":99}`)
			},
			prefix:  "",
			key:     "mykey",
			want:    testStruct{Name: "noprefix", Value: 99},
			wantErr: nil,
		},
		{
			name: "returns ErrCantDecode for invalid JSON",
			setup: func(m *mockStore) {
				m.data["testprefix/badjson"] = []byte(`not valid json`)
			},
			prefix:  "testprefix",
			key:     "badjson",
			wantErr: ErrCantDecode,
			errCheck: func(err error) bool {
				return errors.Is(err, ErrCantDecode)
			},
		},
		{
			name: "unmarshals empty struct",
			setup: func(m *mockStore) {
				m.data["testprefix/empty"] = []byte(`{}`)
			},
			prefix:  "testprefix",
			key:     "empty",
			want:    testStruct{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			j := &JSON[testStruct]{
				Underlying: m,
				Prefix:     tt.prefix,
			}

			got, err := j.Get(context.Background(), tt.key)

			if tt.wantErr != nil {
				if tt.errCheck != nil {
					if !tt.errCheck(err) {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					}
				} else if err != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Get() unexpected error = %v", err)
				}
				if got != tt.want {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestJSON_Set(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		prefix  string
		key     string
		value   testStruct
		setup   func(*mockStore)
		wantErr error
		verify  func(*testing.T, *mockStore)
	}{
		{
			name:    "sets struct with prefix",
			prefix:  "testprefix",
			key:     "mykey",
			value:   testStruct{Name: "test", Value: 42},
			setup:   func(m *mockStore) {},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got, ok := m.data["testprefix/mykey"]
				if !ok {
					t.Error("Set() key was not stored")
				}
				want := []byte(`{"name":"test","value":42}`)
				if !equalBytes(got, want) {
					t.Errorf("Set() stored value = %s, want %s", got, want)
				}
			},
		},
		{
			name:    "sets struct without prefix",
			prefix:  "",
			key:     "mykey",
			value:   testStruct{Name: "noprefix", Value: 99},
			setup:   func(m *mockStore) {},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got, ok := m.data["mykey"]
				if !ok {
					t.Error("Set() key was not stored")
				}
				want := []byte(`{"name":"noprefix","value":99}`)
				if !equalBytes(got, want) {
					t.Errorf("Set() stored value = %s, want %s", got, want)
				}
			},
		},
		{
			name:    "sets empty struct",
			prefix:  "testprefix",
			key:     "empty",
			value:   testStruct{},
			setup:   func(m *mockStore) {},
			wantErr: nil,
			verify: func(t *testing.T, m *mockStore) {
				got, ok := m.data["testprefix/empty"]
				if !ok {
					t.Error("Set() key was not stored")
				}
				want := []byte(`{"name":"","value":0}`)
				if !equalBytes(got, want) {
					t.Errorf("Set() stored value = %s, want %s", got, want)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			j := &JSON[testStruct]{
				Underlying: m,
				Prefix:     tt.prefix,
			}

			err := j.Set(context.Background(), tt.key, tt.value)

			if err != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.verify != nil {
				tt.verify(t, m)
			}
		})
	}
}

func TestJSON_Delete(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		prefix  string
		setup   func(*mockStore)
		key     string
		wantErr error
	}{
		{
			name: "deletes key with prefix",
			setup: func(m *mockStore) {
				m.data["testprefix/mykey"] = []byte(`{"name":"test"}`)
			},
			prefix:  "testprefix",
			key:     "mykey",
			wantErr: nil,
		},
		{
			name:    "deleting non-existent key with prefix is no-op",
			setup:   func(m *mockStore) {},
			prefix:  "testprefix",
			key:     "nonexistent",
			wantErr: nil,
		},
		{
			name: "deletes key without prefix",
			setup: func(m *mockStore) {
				m.data["mykey"] = []byte(`{"name":"test"}`)
			},
			prefix:  "",
			key:     "mykey",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			j := &JSON[testStruct]{
				Underlying: m,
				Prefix:     tt.prefix,
			}

			err := j.Delete(context.Background(), tt.key)

			if err != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify key is gone
			fullKey := tt.key
			if tt.prefix != "" {
				fullKey = tt.prefix + "/" + tt.key
			}
			if _, exists := m.data[fullKey]; exists {
				t.Errorf("Delete() key %q still exists in store", fullKey)
			}
		})
	}
}

func TestJSON_List(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		prefix  string
		listArg string
		setup   func(*mockStore)
		want    []string
		wantErr error
	}{
		{
			name: "lists keys with matching prefix",
			setup: func(m *mockStore) {
				m.data["testprefix/a"] = []byte(`{"name":"a"}`)
				m.data["testprefix/b"] = []byte(`{"name":"b"}`)
				m.data["testprefix/c"] = []byte(`{"name":"c"}`)
				m.data["otherprefix/x"] = []byte(`{"name":"x"}`)
			},
			prefix:  "testprefix",
			listArg: "",
			want:    []string{"a", "b", "c"},
			wantErr: nil,
		},
		{
			name: "lists keys with sub-prefix",
			setup: func(m *mockStore) {
				m.data["testprefix/sub/a"] = []byte(`{"name":"a"}`)
				m.data["testprefix/sub/b"] = []byte(`{"name":"b"}`)
				m.data["testprefix/other/x"] = []byte(`{"name":"x"}`)
			},
			prefix:  "testprefix",
			listArg: "sub/",
			want:    []string{"a", "b"},
			wantErr: nil,
		},
		{
			name:    "returns empty list for non-existent prefix",
			setup:   func(m *mockStore) {},
			prefix:  "testprefix",
			listArg: "noprefix/",
			want:    []string{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMockStore()
			tt.setup(m)

			j := &JSON[testStruct]{
				Underlying: m,
				Prefix:     tt.prefix,
			}

			got, err := j.List(context.Background(), tt.listArg)

			if err != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !equalStringSlicesUnordered(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// equalBytes compares byte slices for equality.
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// equalStringSlicesUnordered compares string slices ignoring order.
func equalStringSlicesUnordered(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]bool)
	for _, s := range a {
		aMap[s] = true
	}
	for _, s := range b {
		if !aMap[s] {
			return false
		}
	}
	return true
}
