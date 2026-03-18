package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container is a simple dependency injection container
type Container struct {
	mu       sync.RWMutex
	services map[string]interface{}
	factories map[string]func() interface{}
	singletons map[string]interface{}
	resolved  map[string]bool
}

// New creates a new container
func New() *Container {
	return &Container{
		services:   make(map[string]interface{}),
		factories:  make(map[string]func() interface{}),
		singletons: make(map[string]interface{}),
		resolved:   make(map[string]bool),
	}
}

// Register registers a service instance
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// RegisterFactory registers a factory function
func (c *Container) RegisterFactory(name string, factory func() interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.factories[name] = factory
}

// RegisterSingleton registers a singleton factory
func (c *Container) RegisterSingleton(name string, factory func() interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.factories[name] = func() interface{} {
		if instance, exists := c.singletons[name]; exists {
			return instance
		}
		instance := factory()
		c.singletons[name] = instance
		return instance
	}
}

// Resolve resolves a service by name
func (c *Container) Resolve(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check if already resolved as singleton
	if instance, exists := c.singletons[name]; exists {
		return instance, nil
	}

	// Check registered services
	if service, exists := c.services[name]; exists {
		return service, nil
	}

	// Check factories
	if factory, exists := c.factories[name]; exists {
		instance := factory()
		return instance, nil
	}

	return nil, fmt.Errorf("service %s not found", name)
}

// MustResolve resolves a service or panics
func (c *Container) MustResolve(name string) interface{} {
	service, err := c.Resolve(name)
	if err != nil {
		panic(err)
	}
	return service
}

// ResolveInto resolves a service into a target
func (c *Container) ResolveInto(name string, target interface{}) error {
	service, err := c.Resolve(name)
	if err != nil {
		return err
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	serviceValue := reflect.ValueOf(service)
	if !serviceValue.Type().AssignableTo(targetValue.Elem().Type()) {
		return fmt.Errorf("service type %s is not assignable to target type %s",
			serviceValue.Type(), targetValue.Elem().Type())
	}

	targetValue.Elem().Set(serviceValue)
	return nil
}

// Has checks if a service is registered
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.services[name]
	if !exists {
		_, exists = c.factories[name]
	}
	return exists
}

// Remove removes a service
func (c *Container) Remove(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.services, name)
	delete(c.factories, name)
	delete(c.singletons, name)
	delete(c.resolved, name)
}

// Clear clears all services
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services = make(map[string]interface{})
	c.factories = make(map[string]func() interface{})
	c.singletons = make(map[string]interface{})
	c.resolved = make(map[string]bool)
}

// List lists all registered services
func (c *Container) List() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make(map[string]bool)
	for name := range c.services {
		names[name] = true
	}
	for name := range c.factories {
		names[name] = true
	}

	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	return result
}

// Service names
const (
	ServiceConfig      = "config"
	ServiceDatabase    = "database"
	ServiceEventBus    = "eventbus"
	ServicePluginMgr   = "plugin_manager"
	ServiceScheduler   = "scheduler"
	ServiceLogger      = "logger"
	ServiceMetrics     = "metrics"
	ServiceAuthService = "auth_service"
	ServiceUserService = "user_service"
	ServiceNodeService = "node_service"
)

// Global container
var globalContainer = New()

// Register registers a service in the global container
func Register(name string, service interface{}) {
	globalContainer.Register(name, service)
}

// RegisterFactory registers a factory in the global container
func RegisterFactory(name string, factory func() interface{}) {
	globalContainer.RegisterFactory(name, factory)
}

// RegisterSingleton registers a singleton in the global container
func RegisterSingleton(name string, factory func() interface{}) {
	globalContainer.RegisterSingleton(name, factory)
}

// Resolve resolves a service from the global container
func Resolve(name string) (interface{}, error) {
	return globalContainer.Resolve(name)
}

// MustResolve resolves a service from the global container or panics
func MustResolve(name string) interface{} {
	return globalContainer.MustResolve(name)
}

// Has checks if a service exists in the global container
func Has(name string) bool {
	return globalContainer.Has(name)
}

// Remove removes a service from the global container
func Remove(name string) {
	globalContainer.Remove(name)
}

// Clear clears the global container
func Clear() {
	globalContainer.Clear()
}

// List lists all services in the global container
func List() []string {
	return globalContainer.List()
}