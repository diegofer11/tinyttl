package tinyttl

import "testing"

func TestNew(t *testing.T) {
	cache := New()
	if cache == nil {
		t.Error("Expected New() to return a non-nil Cache instance")
	}
}
