package config

import "os"

type Config struct {
	DSN            string
	Addr           string
	MigrationsPath string
}

func Load() Config {
	return Config{
		DSN:            get("DB_DSN", "root:root@tcp(127.0.0.1:3307)/article?parseTime=true&multiStatements=true"),
		Addr:           get("ADDR", ":8080"),
		MigrationsPath: get("MIGRATIONS_PATH", "file://migrations"),
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
