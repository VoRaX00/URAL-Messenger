package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upChatsPersons, downChatsPersons)
}

func upChatsPersons(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS chats_persons (
    	chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    	person_id UUID NOT NULL,
    	PRIMARY KEY (chat_id, person_id)
	)`

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downChatsPersons(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS chats_persons`
	_, err := tx.ExecContext(ctx, query)
	return err
}
