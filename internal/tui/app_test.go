package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.currentView != ViewDashboard {
		t.Errorf("expected currentView to be 'dashboard', got %s", m.currentView)
	}

	if len(m.views) != 7 {
		t.Errorf("expected 7 views, got %d", len(m.views))
	}

	if m.showHelp {
		t.Error("showHelp should be false initially")
	}
}

func TestModelInit(t *testing.T) {
	m := NewModel()
	cmd := m.Init()

	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestModelView(t *testing.T) {
	m := NewModel()
	view := m.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	// Check that the view contains expected elements
	if !containsString(view, "AnixOps") {
		t.Error("View should contain 'AnixOps'")
	}
}

func TestModelHelpToggle(t *testing.T) {
	m := NewModel()

	// Toggle help
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model := m2.(Model)

	if !model.showHelp {
		t.Error("showHelp should be true after pressing '?'")
	}
}

func TestModelSearchToggle(t *testing.T) {
	m := NewModel()

	// Toggle search
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model := m2.(Model)

	if !model.showSearch {
		t.Error("showSearch should be true after pressing '/'")
	}
}

func TestModelQuit(t *testing.T) {
	m := NewModel()

	// Press 'q' to quit
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if cmd == nil {
		t.Error("Update should return a command when quitting")
	}
}

func TestModelTabNavigation(t *testing.T) {
	m := NewModel()

	// Press Tab
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := m2.(Model)

	if model.viewIndex != 1 {
		t.Errorf("expected viewIndex 1 after Tab, got %d", model.viewIndex)
	}

	if model.currentView != ViewNodes {
		t.Errorf("expected currentView 'nodes', got %s", model.currentView)
	}
}

func TestModelShiftTabNavigation(t *testing.T) {
	m := NewModel()
	m.viewIndex = 1

	// Press Shift+Tab
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model := m2.(Model)

	if model.viewIndex != 0 {
		t.Errorf("expected viewIndex 0 after Shift+Tab, got %d", model.viewIndex)
	}
}

func TestModelWindowSize(t *testing.T) {
	m := NewModel()

	// Send window size message
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	model := m2.(Model)

	if model.width != 100 {
		t.Errorf("expected width 100, got %d", model.width)
	}

	if model.height != 50 {
		t.Errorf("expected height 50, got %d", model.height)
	}
}

func TestViewConstants(t *testing.T) {
	views := []View{
		ViewDashboard,
		ViewNodes,
		ViewAgents,
		ViewUsers,
		ViewPlaybooks,
		ViewLogs,
		ViewSettings,
	}

	for _, v := range views {
		if string(v) == "" {
			t.Error("View constant should not be empty")
		}
	}
}

func TestNodeItem(t *testing.T) {
	node := NodeItem{
		name:    "test-node",
		host:    "192.168.1.1",
		status:  "online",
		users:   100,
		traffic: "1TB",
	}

	if node.FilterValue() != "test-node" {
		t.Errorf("expected FilterValue 'test-node', got %s", node.FilterValue())
	}

	if node.Title() != "test-node" {
		t.Errorf("expected Title 'test-node', got %s", node.Title())
	}

	if node.Description() == "" {
		t.Error("Description should not be empty")
	}
}

func TestPlaybookItem(t *testing.T) {
	playbook := PlaybookItem{
		name:    "deploy.yml",
		lastRun: "2024-03-15",
		status:  "success",
	}

	if playbook.FilterValue() != "deploy.yml" {
		t.Errorf("expected FilterValue 'deploy.yml', got %s", playbook.FilterValue())
	}

	if playbook.Title() != "deploy.yml" {
		t.Errorf("expected Title 'deploy.yml', got %s", playbook.Title())
	}

	if playbook.Description() == "" {
		t.Error("Description should not be empty")
	}
}

func TestLogEntry(t *testing.T) {
	log := LogEntry{
		Time:    "12:34:56",
		Level:   "INFO",
		Source:  "system",
		Message: "Test message",
	}

	if log.Time != "12:34:56" {
		t.Errorf("expected Time '12:34:56', got %s", log.Time)
	}

	if log.Level != "INFO" {
		t.Errorf("expected Level 'INFO', got %s", log.Level)
	}
}

func TestGenerateSampleLogs(t *testing.T) {
	logs := generateSampleLogs()

	if len(logs) == 0 {
		t.Error("generateSampleLogs should return at least one log entry")
	}

	for _, log := range logs {
		if log.Time == "" || log.Level == "" || log.Message == "" {
			t.Error("LogEntry should have Time, Level, and Message")
		}
	}
}

func TestRenderHelp(t *testing.T) {
	m := NewModel()
	m.showHelp = true

	view := m.View()

	if view == "" {
		t.Error("renderHelp should return non-empty string")
	}

	// Help should contain navigation info
	if !containsString(view, "Navigation") {
		t.Error("Help view should contain 'Navigation'")
	}
}

func TestDefaultKeyMap(t *testing.T) {
	km := defaultKeyMap()

	// Just verify the keymap is not nil and has expected keys
	if km.Quit.Keys() == nil {
		t.Error("Quit key binding should have keys")
	}

	if km.Enter.Keys() == nil {
		t.Error("Enter key binding should have keys")
	}

	if km.Tab.Keys() == nil {
		t.Error("Tab key binding should have keys")
	}
}

func TestShortHelp(t *testing.T) {
	m := NewModel()
	help := m.ShortHelp()

	if len(help) == 0 {
		t.Error("ShortHelp should return at least one binding")
	}
}

func TestFullHelp(t *testing.T) {
	m := NewModel()
	help := m.FullHelp()

	if len(help) == 0 {
		t.Error("FullHelp should return at least one group")
	}
}

func TestModelState(t *testing.T) {
	m := NewModel()

	if m.state.searchInput.Placeholder != "Search..." {
		t.Error("searchInput placeholder should be 'Search...'")
	}
}

func TestViewsList(t *testing.T) {
	m := NewModel()

	expectedViews := []View{
		ViewDashboard, ViewNodes, ViewAgents, ViewUsers,
		ViewPlaybooks, ViewLogs, ViewSettings,
	}

	if len(m.views) != len(expectedViews) {
		t.Errorf("expected %d views, got %d", len(expectedViews), len(m.views))
	}

	for i, v := range expectedViews {
		if m.views[i] != v {
			t.Errorf("expected view %s at index %d, got %s", v, i, m.views[i])
		}
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsString(s[1:], substr) || len(s) >= len(substr) && s[:len(substr)] == substr)
}