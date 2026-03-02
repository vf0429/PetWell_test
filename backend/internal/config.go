package internal

import (
	"os"
	"strings"
)

type AppConfig struct {
	Environment string
	StoreType   string
	Addr        string
	DBPath      string
}

func LoadConfig() AppConfig {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("PETWELL_ENV")))
	if env == "" {
		env = "dev"
	}

	cfg := AppConfig{
		Environment: env,
		StoreType:   "sqlite",
		Addr:        ":8090",
		DBPath:      "./db/petwell_merchant.db",
	}

	switch env {
	case "staging":
		cfg.DBPath = "./db/petwell_merchant_staging.db"
	case "prod":
		cfg.DBPath = "./db/petwell_merchant_prod.db"
	}

	if override := strings.TrimSpace(os.Getenv("PETWELL_STORE")); override != "" {
		cfg.StoreType = strings.ToLower(override)
	}
	if override := strings.TrimSpace(os.Getenv("PETWELL_ADDR")); override != "" {
		cfg.Addr = override
	}
	if override := strings.TrimSpace(os.Getenv("PETWELL_DB_PATH")); override != "" {
		cfg.DBPath = override
	}

	if cfg.StoreType == "memory" {
		cfg.DBPath = ""
	}

	return cfg
}
