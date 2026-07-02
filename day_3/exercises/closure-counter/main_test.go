package main

import "testing"

func TestCounterFun(t *testing.T) {
	// Test 1: Prove that a single counter remembers its state
	t.Run("Single Counter State Retention", func(t *testing.T) {
		counter := counterFun()

		if got := counter(); got != 1 {
			t.Errorf("First call: Expected 1, got %d", got)
		}
		if got := counter(); got != 2 {
			t.Errorf("Second call: Expected 2, got %d", got)
		}
		if got := counter(); got != 3 {
			t.Errorf("Third call: Expected 3, got %d", got)
		}
	})

	// Test 2: Prove that multiple counters do not interfere with each other
	t.Run("Multiple Isolated Counters", func(t *testing.T) {
		counterA := counterFun()
		counterB := counterFun()

		// Advance counterA twice (it should now be at 2)
		counterA() 
		counterA() 

		// counterB is brand new, it MUST start at 1
		if got := counterB(); got != 1 {
			t.Errorf("Expected counterB to start at 1, got %d", got)
		}

		// Advance counterA again, it should remember it was at 2 and go to 3
		if got := counterA(); got != 3 {
			t.Errorf("Expected counterA to continue to 3, got %d", got)
		}
	})
}