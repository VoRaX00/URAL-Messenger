package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upMessages, downMessages)
}

func upMessages(ctx context.Context, tx *sql.Tx) error {
	query := `
	DO $$
		BEGIN 
		    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_status') THEN
				CREATE TYPE message_status as ENUM ('not read', 'read it', 'deleted');
			END IF;
	END $$;
	
	CREATE TABLE IF NOT EXISTS messages (
    	id UUID PRIMARY KEY NOT NULL,
    	message TEXT NOT NULL,
    	person_id UUID NOT NULL,
    	chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    	sending_time TIMESTAMP NOT NULL,
    	status message_status DEFAULT 'not read'
	)`

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downMessages(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS messages`
	_, err := tx.ExecContext(ctx, query)
	return err
}
