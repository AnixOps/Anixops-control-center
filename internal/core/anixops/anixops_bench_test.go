package anixops

import (
	"context"
	"testing"

	"github.com/anixops/anixops-control-center/internal/core/plugin"
)

// BenchmarkMockPlugin is a minimal plugin for benchmarking
type BenchmarkMockPlugin struct{}

func (b *BenchmarkMockPlugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{Name: "bench", Version: "1.0.0"}
}
func (b *BenchmarkMockPlugin) Init(ctx context.Context, config map[string]interface{}) error {
	return nil
}
func (b *BenchmarkMockPlugin) Start(ctx context.Context) error {
	return nil
}
func (b *BenchmarkMockPlugin) Stop(ctx context.Context) error {
	return nil
}
func (b *BenchmarkMockPlugin) HealthCheck(ctx context.Context) error {
	return nil
}
func (b *BenchmarkMockPlugin) Capabilities() []string {
	return nil
}

func BenchmarkNewSDK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = New(nil)
	}
}

func BenchmarkSDKWithConfig(b *testing.B) {
	opts := &Options{
		LogLevel:    "info",
		JSONLogging: false,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New(opts)
	}
}

func BenchmarkRegisterPlugin(b *testing.B) {
	sdk, _ := New(nil)
	p := &BenchmarkMockPlugin{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := string(rune(i))
		_ = sdk.RegisterPlugin(name, p)
	}
}

func BenchmarkStartSDK(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sdk, _ := New(nil)
		sdk.RegisterPlugin("bench", &BenchmarkMockPlugin{})
		_ = sdk.Start(context.Background())
		_ = sdk.Stop(context.Background())
	}
}

func BenchmarkPluginStart(b *testing.B) {
	sdk, _ := New(nil)
	sdk.RegisterPlugin("bench", &BenchmarkMockPlugin{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sdk.Start(context.Background())
		_ = sdk.Stop(context.Background())
	}
}

func BenchmarkHealthCheck(b *testing.B) {
	sdk, _ := New(nil)
	sdk.RegisterPlugin("bench", &BenchmarkMockPlugin{})
	sdk.Start(context.Background())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sdk.HealthCheck(ctx)
	}
}

func BenchmarkGetPluginInfo(b *testing.B) {
	sdk, _ := New(nil)
	sdk.RegisterPlugin("bench", &BenchmarkMockPlugin{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sdk.GetPluginInfo()
	}
}