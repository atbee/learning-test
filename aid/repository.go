package aid

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func StoreAnonymousID(db *sql.DB, tableName string) func(context.Context, time.Time, string, string) error {
	return func(ctx context.Context, now time.Time, did, aid string) error {
		stmt := fmt.Sprintf(`
		INSERT into %s 
		(
			device_id,
			anonymous_id,
			created_at,
		)
		VALUES	($1,$2,$3);
	`, tableName)

		result, err := db.ExecContext(ctx, stmt, did, aid, now)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		_ = rows

		return nil
	}
}
