package goflextime

import (
	_ "embed"
	"fmt"

	"github.com/Songmu/flextime"
	"github.com/jmoiron/sqlx"
)

//go:embed update.sql
var updateSQL string

func UpdateAlreadyRead(tx *sqlx.Tx, userID string) (int, error) {
	now := flextime.Now() // 現在時刻取得

	row := tx.QueryRow(updateSQL, now, userID)

	var updateCnt int64
	if err := row.Scan(&updateCnt); err != nil {
		return 0, fmt.Errorf("update read status: %w", err)
	}

	return int(updateCnt), nil
}
