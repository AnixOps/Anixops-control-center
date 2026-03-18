package main

import (
	"context"
	"embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// Application metadata
const (
	AppName    = "AnixOps Control Center"
	AppVersion = "0.9.9"
	AppID      = "com.anixops.controlcenter"
)

func main() {
	// Create an instance of the app structure
	app := NewApp()
	services := NewAppServices()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     AppName,
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},

		// Performance optimizations
		Frameless: false,
		StartHidden: false,
		HideWindowOnClose: false,

		// Enable debug in development
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},

		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			services.startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			app.shutdown(ctx)
			// services.shutdown(ctx) // TODO: implement shutdown
		},
		Bind: []interface{}{
			app,
			services,
			services.config,
		},

		// Windows specific options
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			Theme:                             windows.Dark,
			// Performance optimizations
			WebviewBrowserPath: "",
		},

		// macOS specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
			},
			Appearance:           mac.DarkAppearance,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   AppName,
				Message: "Unified Control Center for AnixOps Products\n\nVersion " + AppVersion,
			},
			// Performance optimizations
			Preferences: &mac.Preferences{
				MinimumSystemVersion: "10.15",
			},
		},

		// Linux specific options
		Linux: &linux.Options{
			ProgramName: AppName,
			// Performance optimizations
			WindowStartState: linux.Normal,
		},
	})

	if err != nil {
		log.Println("Error:", err.Error())
	}
}

// getPlatformInfo returns information about the current platform
func getPlatformInfo() map[string]string {
	return map[string]string{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"version": AppVersion,
	}
}