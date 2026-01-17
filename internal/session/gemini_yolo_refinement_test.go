package session

import (
	"os/exec"
	"testing"
)

func TestInstance_SetGeminiYoloMode(t *testing.T) {
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not available")
	}

	// Create and start instance
	inst := NewInstanceWithTool("yolo-sync-test", "/tmp", "gemini")
	err := inst.Start()
	if err != nil {
		t.Fatalf("Failed to start instance: %v", err)
	}
	defer func() { _ = inst.Kill() }()

	// Initially YOLO should be nil/false
	if inst.GeminiYoloMode != nil && *inst.GeminiYoloMode {
		t.Error("GeminiYoloMode should not be true initially")
	}

	// Set YOLO mode to true
	inst.SetGeminiYoloMode(true)

	// 1. Verify struct is updated
	if inst.GeminiYoloMode == nil || !*inst.GeminiYoloMode {
		t.Error("SetGeminiYoloMode(true) did not update struct")
	}

	// 2. Verify tmux environment is synced
	tmuxSess := inst.GetTmuxSession()
	val, err := tmuxSess.GetEnvironment("GEMINI_YOLO_MODE")
	if err != nil {
		t.Errorf("Failed to get GEMINI_YOLO_MODE from tmux: %v", err)
	}
	if val != "true" {
		t.Errorf("tmux env GEMINI_YOLO_MODE = %q, want \"true\"", val)
	}

	// Set YOLO mode back to false
	inst.SetGeminiYoloMode(false)
	if inst.GeminiYoloMode == nil || *inst.GeminiYoloMode {
		t.Error("SetGeminiYoloMode(false) did not update struct")
	}
	val, _ = tmuxSess.GetEnvironment("GEMINI_YOLO_MODE")
	if val != "false" {
		t.Errorf("tmux env GEMINI_YOLO_MODE = %q, want \"false\"", val)
	}
}

func TestInstance_UpdateGeminiSession_DetectsYoloFromEnv(t *testing.T) {
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not available")
	}

	inst := NewInstanceWithTool("yolo-detect-test", "/tmp", "gemini")
	err := inst.Start()
	if err != nil {
		t.Fatalf("Failed to start instance: %v", err)
	}
	defer func() { _ = inst.Kill() }()

	tmuxSess := inst.GetTmuxSession()
	
	// Simulate YOLO mode being enabled in the environment (e.g. by a previous session)
	_ = tmuxSess.SetEnvironment("GEMINI_YOLO_MODE", "true")

	// Struct initially has nil
	inst.GeminiYoloMode = nil

	// Trigger update
	inst.UpdateGeminiSession(nil)

	// Should have detected from environment
	if inst.GeminiYoloMode == nil || !*inst.GeminiYoloMode {
		t.Error("UpdateGeminiSession failed to detect YOLO mode from tmux environment")
	}
}
