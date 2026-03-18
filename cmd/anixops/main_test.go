package main

import (
	"testing"
)

func TestVersionInfo(t *testing.T) {
	// Test that version variables are set
	if version == "" {
		t.Log("version is empty (expected for dev builds)")
	}
	if commit == "" {
		t.Log("commit is empty (expected for dev builds)")
	}
	if date == "" {
		t.Log("date is empty (expected for dev builds)")
	}
}

func TestRootCmd(t *testing.T) {
	if rootCmd.Use != "anixops" {
		t.Errorf("expected rootCmd.Use 'anixops', got %s", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}

func TestServerCmd(t *testing.T) {
	if serverCmd.Use != "server" {
		t.Errorf("expected serverCmd.Use 'server', got %s", serverCmd.Use)
	}

	if serverCmd.Short == "" {
		t.Error("serverCmd.Short should not be empty")
	}
}

func TestTUICmd(t *testing.T) {
	if tuiCmd.Use != "tui" {
		t.Errorf("expected tuiCmd.Use 'tui', got %s", tuiCmd.Use)
	}

	if tuiCmd.Short == "" {
		t.Error("tuiCmd.Short should not be empty")
	}
}

func TestAnsibleCmd(t *testing.T) {
	if ansibleCmd.Use != "ansible" {
		t.Errorf("expected ansibleCmd.Use 'ansible', got %s", ansibleCmd.Use)
	}

	if ansibleCmd.Short == "" {
		t.Error("ansibleCmd.Short should not be empty")
	}
}

func TestPlaybookCmd(t *testing.T) {
	if playbookCmd.Use != "run <playbook>" {
		t.Errorf("expected playbookCmd.Use 'run <playbook>', got %s", playbookCmd.Use)
	}

	if playbookCmd.Short == "" {
		t.Error("playbookCmd.Short should not be empty")
	}
}

func TestNodesCmd(t *testing.T) {
	if nodesCmd.Use != "nodes" {
		t.Errorf("expected nodesCmd.Use 'nodes', got %s", nodesCmd.Use)
	}

	if nodesCmd.Short == "" {
		t.Error("nodesCmd.Short should not be empty")
	}
}

func TestListNodesCmd(t *testing.T) {
	if listNodesCmd.Use != "list" {
		t.Errorf("expected listNodesCmd.Use 'list', got %s", listNodesCmd.Use)
	}

	if listNodesCmd.Short == "" {
		t.Error("listNodesCmd.Short should not be empty")
	}
}

func TestVersionCmd(t *testing.T) {
	if versionCmd.Use != "version" {
		t.Errorf("expected versionCmd.Use 'version', got %s", versionCmd.Use)
	}

	if versionCmd.Short == "" {
		t.Error("versionCmd.Short should not be empty")
	}
}

func TestCommandRegistration(t *testing.T) {
	// Check that commands are registered
	commands := rootCmd.Commands()

	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	expectedCommands := []string{"server", "tui", "version", "ansible", "nodes"}
	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("expected command '%s' to be registered", expected)
		}
	}
}

func TestAnsibleSubcommands(t *testing.T) {
	commands := ansibleCmd.Commands()

	if len(commands) == 0 {
		t.Error("ansible command should have subcommands")
	}

	// Check for 'run' subcommand
	found := false
	for _, cmd := range commands {
		if cmd.Name() == "run" {
			found = true
			break
		}
	}

	if !found {
		t.Error("ansible command should have 'run' subcommand")
	}
}

func TestNodesSubcommands(t *testing.T) {
	commands := nodesCmd.Commands()

	if len(commands) == 0 {
		t.Error("nodes command should have subcommands")
	}

	// Check for 'list' subcommand
	found := false
	for _, cmd := range commands {
		if cmd.Name() == "list" {
			found = true
			break
		}
	}

	if !found {
		t.Error("nodes command should have 'list' subcommand")
	}
}

func TestPlaybookFlags(t *testing.T) {
	flags := playbookCmd.Flags()

	// Check that expected flags exist
	expectedFlags := []string{"inventory", "tags", "limit", "check", "verbose"}

	for _, flag := range expectedFlags {
		if flags.Lookup(flag) == nil {
			t.Errorf("expected flag '%s' to be defined", flag)
		}
	}
}

func TestRootFlags(t *testing.T) {
	flags := rootCmd.PersistentFlags()

	if flags.Lookup("config") == nil {
		t.Error("expected 'config' flag to be defined")
	}
}