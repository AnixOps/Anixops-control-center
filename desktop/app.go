package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// App struct
type App struct {
	ctx     context.Context
	version string
	config  *Config
}

// Config holds application configuration
type Config struct {
	APIUrl     string `json:"api_url"`
	Theme      string `json:"theme"`
	AutoUpdate bool   `json:"auto_update"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		version: "1.0.0",
		config: &Config{
			APIUrl:     "http://localhost:8080/api/v1",
			Theme:      "dark",
			AutoUpdate: true,
		},
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Load saved config
	a.loadConfig()
}

// shutdown is called when the app stops
func (a *App) shutdown(ctx context.Context) {
	// Save config
	a.saveConfig()
}

// ============================================
// System Info Methods
// ============================================

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return a.version
}

// GetSystemInfo returns system information
func (a *App) GetSystemInfo() map[string]interface{} {
	return map[string]interface{}{
		"version":    a.version,
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"go_version": runtime.Version(),
		"cpus":       runtime.NumCPU(),
	}
}

// ============================================
// Config Methods
// ============================================

// GetConfig returns the current configuration
func (a *App) GetConfig() *Config {
	return a.config
}

// SetConfig updates the configuration
func (a *App) SetConfig(config Config) {
	a.config = &config
	a.saveConfig()
}

// SetTheme sets the application theme
func (a *App) SetTheme(theme string) {
	a.config.Theme = theme
	a.saveConfig()
}

// SetAPIUrl sets the API URL
func (a *App) SetAPIUrl(url string) {
	a.config.APIUrl = url
	a.saveConfig()
}

func (a *App) loadConfig() {
	// TODO: Load from file
}

func (a *App) saveConfig() {
	// TODO: Save to file
}

// ============================================
// Window Methods
// ============================================

// Minimize minimizes the window
func (a *App) Minimize() {
	// Handled by frontend
}

// Maximize maximizes the window
func (a *App) Maximize() {
	// Handled by frontend
}

// Close closes the application
func (a *App) Close() {
	// Handled by frontend
}

// ============================================
// System Operations
// ============================================

// OpenInBrowser opens a URL in the default browser
func (a *App) OpenInBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// OpenInFileBrowser opens the file browser at a path
func (a *App) OpenInFileBrowser(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// ============================================
// Update Methods
// ============================================

// CheckForUpdates checks for application updates
func (a *App) CheckForUpdates() map[string]interface{} {
	// TODO: Implement update check
	return map[string]interface{}{
		"update_available": false,
		"current_version":  a.version,
		"latest_version":   a.version,
	}
}

// DownloadUpdate downloads the latest update
func (a *App) DownloadUpdate() error {
	// TODO: Implement update download
	return nil
}

// InstallUpdate installs the downloaded update
func (a *App) InstallUpdate() error {
	// TODO: Implement update installation
	return nil
}

// ============================================
// Notification Methods
// ============================================

// SendNotification sends a desktop notification
func (a *App) SendNotification(title, message string) {
	// Handled by frontend using Notification API
}

// ============================================
// File Operations
// ============================================

// SelectFile opens a file selection dialog
func (a *App) SelectFile(title string) (string, error) {
	// TODO: Implement file dialog
	return "", nil
}

// SelectFolder opens a folder selection dialog
func (a *App) SelectFolder(title string) (string, error) {
	// TODO: Implement folder dialog
	return "", nil
}

// SaveFile opens a save file dialog
func (a *App) SaveFile(title, defaultName string) (string, error) {
	// TODO: Implement save dialog
	return "", nil
}

// ============================================
// Logging Methods
// ============================================

// Log writes a log message
func (a *App) Log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] %s\n", timestamp, level, message)
}

// LogInfo writes an info log message
func (a *App) LogInfo(message string) {
	a.Log("INFO", message)
}

// LogError writes an error log message
func (a *App) LogError(message string) {
	a.Log("ERROR", message)
}

// LogDebug writes a debug log message
func (a *App) LogDebug(message string) {
	if os.Getenv("DEBUG") != "" {
		a.Log("DEBUG", message)
	}
}

// ============================================
// Environment Methods
// ============================================

// GetEnv gets an environment variable
func (a *App) GetEnv(key string) string {
	return os.Getenv(key)
}

// SetEnv sets an environment variable
func (a *App) SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// GetHomeDir gets the user's home directory
func (a *App) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// GetCurrentDir gets the current working directory
func (a *App) GetCurrentDir() (string, error) {
	return os.Getwd()
}