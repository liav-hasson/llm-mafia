package names

import (
	"testing"
)

func TestNewGenerator(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie"}
	gen, err := NewGenerator(names)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
}

func TestNewGenerator_EmptyList(t *testing.T) {
	_, err := NewGenerator([]string{})

	if err == nil {
		t.Fatal("expected error for empty names list")
	}
}

func TestNext(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie"}
	gen, _ := NewGenerator(names)

	// Get all names in sequence
	for i, expected := range names {
		name, err := gen.Next()
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if name != expected {
			t.Errorf("at index %d: got %q, expected %q", i, name, expected)
		}
	}

	// Fourth call should return error
	_, err := gen.Next()
	if err != ErrNoMoreNames {
		t.Errorf("expected ErrNoMoreNames, got %v", err)
	}
}

func TestReset(t *testing.T) {
	names := []string{"Alice", "Bob"}
	gen, _ := NewGenerator(names)

	// Use all names
	if _, err := gen.Next(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if _, err := gen.Next(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Reset
	gen.Reset()

	// Should be able to get first name again
	name, err := gen.Next()
	if err != nil {
		t.Fatalf("unexpected error after reset: %v", err)
	}
	if name != "Alice" {
		t.Errorf("after reset: got %q, expected %q", name, "Alice")
	}
}

func TestRemaining(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie"}
	gen, _ := NewGenerator(names)

	if rem := gen.Remaining(); rem != 3 {
		t.Errorf("initial remaining: got %d, expected 3", rem)
	}

	if _, err := gen.Next(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if rem := gen.Remaining(); rem != 2 {
		t.Errorf("after one Next: got %d, expected 2", rem)
	}

	if _, err := gen.Next(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if _, err := gen.Next(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if rem := gen.Remaining(); rem != 0 {
		t.Errorf("after all Next: got %d, expected 0", rem)
	}
}

func TestConcurrency(t *testing.T) {
	names := make([]string, 100)
	for i := 0; i < 100; i++ {
		names[i] = "Player"
	}

	gen, _ := NewGenerator(names)

	// Spawn multiple goroutines calling Next concurrently
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			_, _ = gen.Next()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}

	// Should have consumed all names
	if rem := gen.Remaining(); rem != 0 {
		t.Errorf("after concurrent access: expected 0 remaining, got %d", rem)
	}
}
