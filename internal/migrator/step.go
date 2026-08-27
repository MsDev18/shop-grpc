package migrator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
)

func (m Migrator) Step(step int) error {
	err := m.migrate.Steps(step)
	if err != nil {
		_, isShortLimitErr := err.(migrate.ErrShortLimit)
		switch {
		case errors.Is(err ,migrate.ErrNoChange):
			fmt.Println("no change in migration")
			return nil
		// type assertion for check error type
		case isShortLimitErr:
			fmt.Println("limit reached")
			return nil
		case strings.Contains(err.Error(), "file does not exist"):
			fmt.Println("no more migrations to apply")
			return nil
		default:
			return fmt.Errorf("error in step migration: %w", err)
		}
	}
	fmt.Println("migration successfully")
	return nil
}
