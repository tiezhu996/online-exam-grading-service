package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("EXAM_WORKERS", "")
	t.Setenv("EXAM_BATCH_SIZE", "")
	c := Load()
	if c.Workers != 2 || c.BatchSize != 2 {
		t.Fatalf("bad config %+v", c)
	}
}
