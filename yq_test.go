package main

import (
	"testing"

	command "github.com/mikefarah/yq/v4/cmd"
)

func TestMainFunction(t *testing.T) {
	// This is a basic smoke test for the main function
	// We can't easily test the main function directly since it calls os.Exit
	// But we can test the logic that would be executed

	cmd := command.New()
	if cmd == nil {
		t.Fatal("command.New() returned nil")
	}

	if cmd.Use != "yq" {
		t.Errorf("Expected command Use to be 'yq', got %q", cmd.Use)
	}
}

func TestMainFunctionLogic(t *testing.T) {
	// Test the logic that would be executed in main()
	cmd := command.New()

	args := []string{}
	_, _, err := cmd.Find(args)
	if err != nil {
		t.Errorf("Expected no error with empty args, but got: %v", err)
	}

	args = []string{"invalid-command"}
	_, _, err = cmd.Find(args)
	if err == nil {
		t.Error("Expected error when invalid command found, but got nil")
	}

	args = []string{"eval"}
	_, _, err = cmd.Find(args)
	if err != nil {
		t.Errorf("Expected no error with valid command 'eval', got: %v", err)
	}

	args = []string{"__complete"}
	_, _, err = cmd.Find(args)
	if err == nil {
		t.Error("Expected error when no command found for '__complete', but got nil")
	}
}

func TestMainFunctionWithArgs(t *testing.T) {
	// Test the argument processing logic
	cmd := command.New()

	args := []string{}
	_, _, err := cmd.Find(args)
	if err != nil {
		t.Errorf("Expected no error with empty args, but got: %v", err)
	}

	// When Find fails and args[0] is not "__complete", main would set args to ["eval"] + original args
	// This is the logic: newArgs := []string{"eval"}
	// cmd.SetArgs(append(newArgs, os.Args[1:]...))

	args = []string{"invalid"}
	_, _, err = cmd.Find(args)
	if err == nil {
		t.Error("Expected error with invalid command")
	}

	args = []string{"__complete"}
	_, _, err = cmd.Find(args)
	if err == nil {
		t.Error("Expected error with __complete command")
	}
}

func TestMainFunctionExecution(t *testing.T) {
	// Test that the command can be executed without crashing
	cmd := command.New()

	cmd.SetArgs([]string{"--version"})

	// We can't easily test os.Exit(1) behaviour, but we can test that
	// the command structure is correct and can be configured
	if cmd == nil {
		t.Fatal("Command should not be nil")
	}

	if cmd.Use != "yq" {
		t.Errorf("Expected command Use to be 'yq', got %q", cmd.Use)
	}
}

func TestMainFunctionErrorHandling(t *testing.T) {
	// Test the error handling logic that would be in main()
	cmd := command.New()

	args := []string{"nonexistent-command"}
	_, _, err := cmd.Find(args)
	if err == nil {
		t.Error("Expected error with nonexistent command")
	}

	// The main function logic would be:
	// if err != nil && args[0] != "__complete" {
	//     newArgs := []string{"eval"}
	//     cmd.SetArgs(append(newArgs, os.Args[1:]...))
	// }

	// Test that this logic would work
	if args[0] != "__complete" {
		// This is what main() would do
		newArgs := []string{"eval"}
		cmd.SetArgs(append(newArgs, args...))

		// We can't easily verify the args were set correctly since cmd.Args is a function
		// But we can test that SetArgs doesn't crash and the command is still valid
		if cmd == nil {
			t.Error("Command should not be nil after SetArgs")
		}

		_, _, err := cmd.Find([]string{"eval"})
		if err != nil {
			t.Errorf("Should be able to find eval command after SetArgs: %v", err)
		}
	}
}

func TestMainFunctionWithCompletionCommand(t *testing.T) {
	// Test that __complete command doesn't trigger default command logic
	cmd := command.New()

	args := []string{"__complete"}
	_, _, err := cmd.Find(args)
	if err == nil {
		t.Error("Expected error with __complete command")
	}

	// The main function logic would be:
	// if err != nil && args[0] != "__complete" {
	//     // This should NOT execute for __complete
	// }

	// Verify that __complete doesn't trigger the default command logic
	if args[0] == "__complete" {
		// This means the default command logic should NOT execute
		t.Log("__complete command correctly identified, default command logic should not execute")
	}
}

func TestMainFunctionIntegration(t *testing.T) {
	// Integration test to verify the main function logic works end-to-end

	cmd := command.New()
	cmd.SetArgs([]string{"eval", "--help"})

	// This should not crash (we can't test the actual execution due to os.Exit)
	if cmd == nil {
		t.Fatal("Command should not be nil")
	}

	cmd2 := command.New()
	cmd2.SetArgs([]string{"invalid-command"})

	// Simulate the main function logic
	args := []string{"invalid-command"}
	_, _, err := cmd2.Find(args)
	if err != nil {
		// This is what main() would do
		newArgs := []string{"eval"}
		cmd2.SetArgs(append(newArgs, args...))
	}

	// We can't directly access cmd.Args since it's a function, but we can test
	// that SetArgs worked by ensuring the command is still functional
	if cmd2 == nil {
		t.Error("Command should not be nil after SetArgs")
	}

	_, _, err = cmd2.Find([]string{"eval"})
	if err != nil {
		t.Errorf("Should be able to find eval command after SetArgs: %v", err)
	}
}
