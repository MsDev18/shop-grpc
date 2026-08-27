package migrator

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func (m Migrator) Up() error {
	err := m.migrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("error in up migration : %w", err)
	}
	fmt.Println("Migration up completed successfully")
	return nil
}
