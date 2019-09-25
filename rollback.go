package pqutil

import (
	"database/sql"
	"fmt"
)

// HandleRollback takes an sql.Tx and an error that occured and
// attempts to roll it back. If the rollback fails, the error is
// wrapped and returned.
func HandleRollback(tx *sql.Tx, err error) error {
	rollErr := tx.Rollback()
	if rollErr != nil {
		return fmt.Errorf("rollback failed: %s, %w", rollErr, err)
	}
	return err
}
