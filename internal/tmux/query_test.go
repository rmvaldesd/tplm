package tmux

import "testing"

func TestNeighborSession(t *testing.T) {
	// NeighborSession depends on ListSessions which calls tmux.
	// We test the logic by verifying behavior when tmux is not running
	// (returns no sessions).
	t.Run("no sessions returns false", func(t *testing.T) {
		name, ok := NeighborSession("nonexistent")
		if ok {
			t.Errorf("NeighborSession() ok = true, want false")
		}
		if name != "" {
			t.Errorf("NeighborSession() name = %q, want empty", name)
		}
	})
}

func TestSessionExists(t *testing.T) {
	t.Run("nonexistent session returns false", func(t *testing.T) {
		// This will fail gracefully when tmux is not running.
		if SessionExists("__tplm_test_nonexistent__") {
			t.Error("SessionExists() = true for nonexistent session")
		}
	})
}
