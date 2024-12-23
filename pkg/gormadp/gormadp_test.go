package gormadp

import (
	"testing"
)

func Test_New(t *testing.T) {
	cfg := ConfigDefault
	options := cfg.Dialect.PgOptions(PgConfig{
		SSLMode:  false,
		TimeZone: "Asia/Manila",
	})

	if db := Open(cfg, options); db == nil {
		t.Errorf("dbadapter.Open() go error '%v', expected '%v'", "db = nil", "db *gorm.DB && db != nil")
	}
}

func Benchmark_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg := ConfigDefault
		options := cfg.Dialect.PgOptions(PgConfig{
			SSLMode:  false,
			TimeZone: "Asia/Manila",
		})

		db := Open(cfg, options)
		if db == nil {
			b.Errorf("dbadapter.Open() go error '%v', expected '%v'", "db = nil", "db *gorm.DB && db != nil")
		}

		db.Close()
	}
}
