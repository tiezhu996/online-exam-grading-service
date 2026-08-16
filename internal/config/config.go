package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers   int
	BatchSize int
}

func Load() *Config {
	return &Config{
		Workers:   getInt("EXAM_WORKERS", 2),
		BatchSize: getInt("EXAM_BATCH_SIZE", 2),
	}
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
