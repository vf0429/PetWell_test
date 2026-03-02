package internal

import "testing"

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("PETWELL_ENV", "")
	t.Setenv("PETWELL_STORE", "")
	t.Setenv("PETWELL_ADDR", "")
	t.Setenv("PETWELL_DB_PATH", "")

	cfg := LoadConfig()
	if cfg.Environment != "dev" {
		t.Fatalf("expected environment dev, got %q", cfg.Environment)
	}
	if cfg.StoreType != "sqlite" {
		t.Fatalf("expected store sqlite, got %q", cfg.StoreType)
	}
	if cfg.Addr != ":8090" {
		t.Fatalf("expected addr :8090, got %q", cfg.Addr)
	}
	if cfg.DBPath != "./db/petwell_merchant.db" {
		t.Fatalf("expected default db path, got %q", cfg.DBPath)
	}
}

func TestLoadConfigProfileDefaults(t *testing.T) {
	t.Setenv("PETWELL_ENV", "staging")
	t.Setenv("PETWELL_STORE", "")
	t.Setenv("PETWELL_ADDR", "")
	t.Setenv("PETWELL_DB_PATH", "")

	cfg := LoadConfig()
	if cfg.Environment != "staging" {
		t.Fatalf("expected environment staging, got %q", cfg.Environment)
	}
	if cfg.DBPath != "./db/petwell_merchant_staging.db" {
		t.Fatalf("expected staging db path, got %q", cfg.DBPath)
	}
}

func TestLoadConfigOverrides(t *testing.T) {
	t.Setenv("PETWELL_ENV", "prod")
	t.Setenv("PETWELL_STORE", "memory")
	t.Setenv("PETWELL_ADDR", ":9090")
	t.Setenv("PETWELL_DB_PATH", "/tmp/custom.db")

	cfg := LoadConfig()
	if cfg.Environment != "prod" {
		t.Fatalf("expected environment prod, got %q", cfg.Environment)
	}
	if cfg.StoreType != "memory" {
		t.Fatalf("expected store memory, got %q", cfg.StoreType)
	}
	if cfg.Addr != ":9090" {
		t.Fatalf("expected addr :9090, got %q", cfg.Addr)
	}
	if cfg.DBPath != "" {
		t.Fatalf("expected empty db path for memory store, got %q", cfg.DBPath)
	}
}
