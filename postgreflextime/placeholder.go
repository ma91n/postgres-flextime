package example

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//go:embed update.sql
var updateSQL string

func UpdateAlreadyRead(tx *sqlx.Tx, userID string) (int, error) {
	row := tx.QueryRow(updateSQL, userID)

	var updateCnt int64
	if err := row.Scan(&updateCnt); err != nil {
		return 0, fmt.Errorf("update read status: %w", err)
	}

	return int(updateCnt), nil
}
