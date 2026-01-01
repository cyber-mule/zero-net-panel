package config

import (
	"testing"
)

func TestMetricsConfigNormalize(t *testing.T) {
	cfg := MetricsConfig{Enable: true, Path: "metrics", ListenOn: " 0.0.0.0:9100 "}
	cfg.Normalize()

	if cfg.Path != "/metrics" {
		t.Fatalf("expected /metrics path, got %q", cfg.Path)
	}
	if cfg.ListenOn != "0.0.0.0:9100" {
		t.Fatalf("expected trimmed listen address, got %q", cfg.ListenOn)
	}
}

func TestMetricsConfigNormalizeDisabled(t *testing.T) {
	cfg := MetricsConfig{Enable: false, Path: "  /custom ", ListenOn: "127.0.0.1:9000"}
	cfg.Normalize()

	if cfg.Path != "/custom" {
		t.Fatalf("expected sanitized path, got %q", cfg.Path)
	}
	if cfg.ListenOn != "" {
		t.Fatalf("expected listen address cleared when disabled, got %q", cfg.ListenOn)
	}
}

func TestConfigNormalizeSyncsMiddlewares(t *testing.T) {
	cfg := Config{}
	cfg.Metrics.Enable = true
	cfg.Metrics.Path = "metrics"

	cfg.Normalize()

	if !cfg.Middlewares.Prometheus {
		t.Fatal("prometheus middleware should be enabled when metrics is on")
	}
	if !cfg.Middlewares.Metrics {
		t.Fatal("metrics middleware should be enabled when metrics is on")
	}

	cfg.Metrics.Enable = false
	cfg.Normalize()

	if cfg.Middlewares.Prometheus {
		t.Fatal("prometheus middleware should be disabled when metrics is off")
	}
	if cfg.Middlewares.Metrics {
		t.Fatal("metrics middleware should be disabled when metrics is off")
	}
}

func TestConfigNormalizeSiteDefaults(t *testing.T) {
	cfg := Config{
		Project: ProjectConfig{
			Name: " Zero Network Panel ",
		},
		Site: SiteConfig{
			Name:    "   ",
			LogoURL: " https://example.com/logo.png ",
		},
	}

	cfg.Normalize()

	if cfg.Site.Name != "Zero Network Panel" {
		t.Fatalf("expected site name to default to project name, got %q", cfg.Site.Name)
	}
	if cfg.Site.LogoURL != "https://example.com/logo.png" {
		t.Fatalf("expected trimmed logo url, got %q", cfg.Site.LogoURL)
	}
}

func TestCORSConfigNormalize(t *testing.T) {
	cfg := CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{" http://localhost:5173 ", "", "http://localhost:5173"},
		AllowHeaders: []string{" X-ZNP-API-Key ", "x-znp-api-key", "Authorization"},
	}
	cfg.Normalize()

	if len(cfg.AllowOrigins) != 1 || cfg.AllowOrigins[0] != "http://localhost:5173" {
		t.Fatalf("expected trimmed origin list, got %#v", cfg.AllowOrigins)
	}
	if len(cfg.AllowHeaders) != 2 {
		t.Fatalf("expected de-duplicated headers, got %#v", cfg.AllowHeaders)
	}
}
