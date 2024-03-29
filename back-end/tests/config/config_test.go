package config_test

import (
	"testing"

	"github.com/andrii-stp/users-crud/config"
)

func TestLoad(t *testing.T) {
	configFile := "../test.env"

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Port expected as 8080 but got %v", cfg.Server.Port)
	}

	if cfg.Database.Driver != "postgres" {
		t.Errorf("Driver expected as postgres but got %v", cfg.Database.Driver)
	}

	if cfg.Database.User != "postgres" {
		t.Errorf("User expected as postgres but got %v", cfg.Database.User)
	}

	if cfg.Database.Password != "postgres" {
		t.Errorf("Password expected as postgres but got %v", cfg.Database.Password)
	}

	if cfg.Database.Name != "users_test" {
		t.Errorf("Name name expected as users_test but got %v", cfg.Database.Name)
	}

	if cfg.Database.SSLMode != "disable" {
		t.Errorf("SSLMode expected as postgres but got %v", cfg.Database.SSLMode)
	}
}
