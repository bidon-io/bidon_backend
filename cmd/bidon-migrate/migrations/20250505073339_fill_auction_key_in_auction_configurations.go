package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"math/big"
	"strings"
)

func init() {
	goose.AddMigrationContext(upFillAuctionKeyInAuctionConfigurations, downFillAuctionKeyInAuctionConfigurations)
}

func upFillAuctionKeyInAuctionConfigurations(ctx context.Context, tx *sql.Tx) error {
	const batchSize = 1000

	for {
		rows, err := tx.Query(`
            SELECT id, public_uid
            FROM auction_configurations
            WHERE public_uid IS NOT NULL
              AND auction_key IS NULL
            ORDER BY id
            LIMIT $1
        `, batchSize)
		if err != nil {
			return fmt.Errorf("select batch: %w", err)
		}

		var (
			ids       []int64
			keys      []string
			id        int64
			publicUID sql.NullInt64
		)
		for rows.Next() {
			if err := rows.Scan(&id, &publicUID); err != nil {
				rows.Close()
				return fmt.Errorf("scan row: %w", err)
			}

			bi := big.NewInt(publicUID.Int64)
			key := strings.ToUpper(bi.Text(32))

			ids = append(ids, id)
			keys = append(keys, key)
		}
		rows.Close()
		if len(ids) == 0 {
			break
		}

		_, err = tx.Exec(`
            UPDATE auction_configurations AS a
            SET auction_key = u.key
            FROM (
                SELECT UNNEST($1::bigint[]) AS id,
                       UNNEST($2::text[])   AS key
            ) AS u
            WHERE a.id = u.id
        `, pq.Array(ids), pq.Array(keys))
		if err != nil {
			return fmt.Errorf("update batch: %w", err)
		}
	}

	return nil
}

func downFillAuctionKeyInAuctionConfigurations(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
        UPDATE auction_configurations
        SET auction_key = NULL
    `)
	return err
}
