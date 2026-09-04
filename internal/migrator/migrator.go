package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


type Migrator struct {
	migrate *migrate.Migrate
}

func New(dsn string) Migrator {
	newDSN := fmt.Sprintf("mysql://%s&x-migrations-table=migrations", dsn)
	migration, err := migrate.New(
		"file://migrations" , 
		newDSN,
	)
	if err != nil {
		panic(fmt.Errorf("error in migration : %w", err))
	}

	return Migrator{
		migrate: migration,
	}
}

// internal/migrator/migrator.go — یک تابع جدید کنار New اضافه کن
func NewWithSourcePath(sourcePath string, dsn string) Migrator {
	newDSN := fmt.Sprintf("mysql://%s&x-migrations-table=migrations", dsn)
	migration, err := migrate.New(
		fmt.Sprintf("file://%s", sourcePath),
		newDSN,
	)
	if err != nil {
		panic(fmt.Errorf("error in migration : %w", err))
	}

	return Migrator{
		migrate: migration,
	}
}
