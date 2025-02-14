package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFiles, downFiles)
}

func upFiles(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS files (
    	id UUID PRIMARY KEY NOT NULL,
    	file BYTEA NOT NULL,
    	message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE
	)`

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downFiles(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS files`
	_, err := tx.ExecContext(ctx, query)
	return err
}
