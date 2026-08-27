package migrator

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func (m Migrator) Down() error {
	err := m.migrate.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to roll back")
			return nil
		}
		return fmt.Errorf("error in down migration : %w", err)
	}
	fmt.Println("Migration down completed successfully")
	return nil
}
