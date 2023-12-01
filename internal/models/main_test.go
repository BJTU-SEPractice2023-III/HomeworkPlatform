package models

import (
	"homework_platform/internal/bootstrap"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	bootstrap.Sqlite = true
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if res := DB.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		panic(res.Error)
	}
	bootstrap.Sqlite = true
	if err != nil {
		panic(err)
	}

	Migrate()

	m.Run()
}