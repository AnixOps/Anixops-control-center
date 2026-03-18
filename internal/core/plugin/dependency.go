package plugin

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/anixops/anixops-control-center/internal/core/errors"
	"github.com/anixops/anixops-control-center/internal/core/logger"
)

// Dependency represents a plugin dependency
type Dependency struct {
	Name     string `json:"name"`
	Version  string `json:"version"`  // semver constraint, e.g., ">=1.0.0,<2.0.0"
	Required bool   `json:"required"` // if false, plugin can work without it
}

// DependencyManager manages plugin dependencies
type DependencyManager struct {
	mu          sync.RWMutex
	dependencies map[string][]Dependency // plugin name -> dependencies
	providers    map[string]string       // capability -> plugin name
}

// NewDependencyManager creates a new dependency manager
func NewDependencyManager() *DependencyManager {
	return &DependencyManager{
		dependencies: make(map[string][]Dependency),
		providers:    make(map[string]string),
	}
}

// Register registers a plugin's dependencies
func (dm *DependencyManager) Register(pluginName string, deps []Dependency) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.dependencies[pluginName] = deps
}

// RegisterProvider registers a plugin as a capability provider
func (dm *DependencyManager) RegisterProvider(capability, pluginName string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.providers[capability] = pluginName
}

// Resolve resolves dependencies and returns start order
func (dm *DependencyManager) Resolve(plugins []string) ([]string, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	// Build dependency graph
	graph := make(map[string][]string)
	for _, plugin := range plugins {
		deps, exists := dm.dependencies[plugin]
		if !exists {
			graph[plugin] = nil
			continue
		}

		var depNames []string
		for _, dep := range deps {
			// Check if dependency exists
			found := false
			for _, p := range plugins {
				if p == dep.Name {
					depNames = append(depNames, dep.Name)
					found = true
					break
				}
			}
			if !found && dep.Required {
				return nil, errors.NewError(errors.CodePluginDependency,
					fmt.Sprintf("required dependency %s not found for plugin %s", dep.Name, plugin),
					errors.LevelError, 500)
			}
		}
		graph[plugin] = depNames
	}

	// Topological sort
	return topologicalSort(graph)
}

// GetProvider returns the plugin that provides a capability
func (dm *DependencyManager) GetProvider(capability string) (string, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	name, ok := dm.providers[capability]
	return name, ok
}

// CheckVersion checks if a version satisfies a constraint
func CheckVersion(version, constraint string) bool {
	// Simple version check - in production, use proper semver library
	if constraint == "" || constraint == "*" {
		return true
	}

	// Simple prefix match for now
	if strings.HasPrefix(version, constraint) {
		return true
	}

	// Handle operators
	if strings.HasPrefix(constraint, ">=") {
		return version >= constraint[2:]
	}
	if strings.HasPrefix(constraint, ">") {
		return version > constraint[1:]
	}
	if strings.HasPrefix(constraint, "<=") {
		return version <= constraint[2:]
	}
	if strings.HasPrefix(constraint, "<") {
		return version < constraint[1:]
	}
	if strings.HasPrefix(constraint, "==") || strings.HasPrefix(constraint, "=") {
		vers := strings.TrimPrefix(strings.TrimPrefix(constraint, "=="), "=")
		return version == vers
	}

	return version == constraint
}

// ValidateDependencies validates all dependencies
func (dm *DependencyManager) ValidateDependencies(plugins map[string]Plugin) []error {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	var errs []error
	for name, p := range plugins {
		deps, exists := dm.dependencies[name]
		if !exists {
			continue
		}

		for _, dep := range deps {
			depPlugin, exists := plugins[dep.Name]
			if !exists {
				if dep.Required {
					errs = append(errs, fmt.Errorf("plugin %s: missing required dependency %s", name, dep.Name))
				}
				continue
			}

			// Check version
			info := depPlugin.Info()
			if !CheckVersion(info.Version, dep.Version) {
				errs = append(errs, fmt.Errorf("plugin %s: dependency %s version %s does not satisfy constraint %s",
					name, dep.Name, info.Version, dep.Version))
			}

			// Use p to avoid unused variable warning
			_ = p
		}
	}
	return errs
}

// topologicalSort performs topological sort on dependency graph
func topologicalSort(graph map[string][]string) ([]string, error) {
	// Calculate in-degrees
	inDegree := make(map[string]int)
	for node := range graph {
		inDegree[node] = 0
	}
	for _, deps := range graph {
		for _, dep := range deps {
			inDegree[dep]++ // dep is depended upon
		}
	}

	// Find nodes with no dependencies
	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}
	sort.Strings(queue) // Deterministic ordering

	var result []string
	for len(queue) > 0 {
		// Pop from queue
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// Remove edges
		for _, dep := range graph[node] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
				sort.Strings(queue) // Keep sorted
			}
		}
	}

	// Check for cycles
	if len(result) != len(graph) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	// Reverse to get correct order (dependencies first)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// PluginWithDeps is a plugin with dependency support
type PluginWithDeps struct {
	Plugin
	dependencies []Dependency
}

// NewPluginWithDeps creates a plugin with dependencies
func NewPluginWithDeps(p Plugin, deps []Dependency) *PluginWithDeps {
	return &PluginWithDeps{
		Plugin:      p,
		dependencies: deps,
	}
}

// Dependencies returns the plugin's dependencies
func (p *PluginWithDeps) Dependencies() []Dependency {
	return p.dependencies
}

// DependentPlugin is a plugin that can declare dependencies
type DependentPlugin interface {
	Plugin
	Dependencies() []Dependency
}

// HookProvider is a plugin that provides hooks
type HookProvider interface {
	Plugin
	Hooks() []Hook
}

// Hook represents a plugin hook
type Hook struct {
	Name     string
	Priority int
	Type     HookType
	Handler  func(ctx context.Context, data interface{}) error
}

// HookType represents the type of hook
type HookType int

const (
	HookPreInit HookType = iota
	HookPostInit
	HookPreStart
	HookPostStart
	HookPreStop
	HookPostStop
	HookPreExecute
	HookPostExecute
	HookHealthCheck
)

// HookRegistry manages plugin hooks
type HookRegistry struct {
	mu    sync.RWMutex
	hooks map[HookType][]Hook
}

// NewHookRegistry creates a new hook registry
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[HookType][]Hook),
	}
}

// Register registers a hook
func (r *HookRegistry) Register(hook Hook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hooks[hook.Type] = append(r.hooks[hook.Type], hook)

	// Sort by priority
	sort.Slice(r.hooks[hook.Type], func(i, j int) bool {
		return r.hooks[hook.Type][i].Priority < r.hooks[hook.Type][j].Priority
	})
}

// Execute executes hooks of a specific type
func (r *HookRegistry) Execute(ctx context.Context, hookType HookType, data interface{}) error {
	r.mu.RLock()
	hooks := make([]Hook, len(r.hooks[hookType]))
	copy(hooks, r.hooks[hookType])
	r.mu.RUnlock()

	for _, hook := range hooks {
		logger.Debug("Executing hook", logger.F("hook", hook.Name), logger.F("type", hookType))
		if err := hook.Handler(ctx, data); err != nil {
			return errors.NewError(errors.CodeInternal,
				fmt.Sprintf("hook %s failed", hook.Name), errors.LevelError, 500).
				WithCause(err)
		}
	}
	return nil
}

// HookString returns a string representation of a hook type
func (t HookType) String() string {
	switch t {
	case HookPreInit:
		return "pre_init"
	case HookPostInit:
		return "post_init"
	case HookPreStart:
		return "pre_start"
	case HookPostStart:
		return "post_start"
	case HookPreStop:
		return "pre_stop"
	case HookPostStop:
		return "post_stop"
	case HookPreExecute:
		return "pre_execute"
	case HookPostExecute:
		return "post_execute"
	case HookHealthCheck:
		return "health_check"
	default:
		return "unknown"
	}
}