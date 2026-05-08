package service

import (
	"testing"

	"sb-sync/pkg/config"
)

func TestGetServiceConfig(t *testing.T) {
	cfg := GetServiceConfig()

	if cfg.Name != config.ServiceName {
		t.Errorf("Expected service name %s, got %s", config.ServiceName, cfg.Name)
	}
	if cfg.DisplayName != config.ServiceDisplayName {
		t.Errorf("Expected display name %s, got %s", config.ServiceDisplayName, cfg.DisplayName)
	}
	if cfg.Description != config.ServiceDescription {
		t.Errorf("Expected description %s, got %s", config.ServiceDescription, cfg.Description)
	}
}

func TestGetSingBoxBinaryPath(t *testing.T) {
	config.AppConfig.InstallDir = "/test/bin"
	path := GetSingBoxBinaryPath()
	if path != "/test/bin/sing-box" {
		t.Errorf("Expected /test/bin/sing-box, got %s", path)
	}
}

func TestServiceInterface(t *testing.T) {
	var s ServiceInterface
	s = &serviceWrapper{}

	if s == nil {
		t.Error("serviceWrapper should implement ServiceInterface")
	}
}
