package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upChats, downChats)
}

func upChats(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS chats (
    	id UUID PRIMARY KEY NOT NULL,
    	name VARCHAR(40) NOT NULL
	)`

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downChats(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS chats`
	_, err := tx.ExecContext(ctx, query)
	return err
}
