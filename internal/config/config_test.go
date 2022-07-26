package config

import "testing"

func TestConfig_LoadConfig(t *testing.T) {
	c := Config{}

	err := c.LoadConfig("test.yaml")

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if c.Database.Host != "localhost" {
		t.Fatalf("Expected localhost, got %s", c.Database.Host)
	}

	if c.Database.Port != 5432 {
		t.Fatalf("Expected 5432, got %d", c.Database.Port)
	}

	if c.Server.Host != "0.0.0.0" {
		t.Fatalf("Expected 0.0.0.0, got %s", c.Server.Host)
	}

	if c.Server.Port != 8090 {
		t.Fatalf("Expected 8090, got %d", c.Server.Port)
	}

	if len(c.Server.AllowedOrigins) != 2 {
		t.Fatalf("Expected 2 allowed origins, got %d", len(c.Server.AllowedOrigins))
	}
}
