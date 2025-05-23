// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: bills.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const getBills = `-- name: GetBills :many

SELECT amount, name, tmstmp
FROM bills
WHERE user_id = $1
AND category = $2
`

type GetBillsParams struct {
	UserID   uuid.UUID
	Category string
}

type GetBillsRow struct {
	Amount float64
	Name   string
	Tmstmp time.Time
}

// AND tmstmp >= $2
// AND tmstmp < $3;
func (q *Queries) GetBills(ctx context.Context, arg GetBillsParams) ([]GetBillsRow, error) {
	rows, err := q.db.QueryContext(ctx, getBills, arg.UserID, arg.Category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBillsRow
	for rows.Next() {
		var i GetBillsRow
		if err := rows.Scan(&i.Amount, &i.Name, &i.Tmstmp); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReport = `-- name: GetReport :many
SELECT category, amount
FROM bills
WHERE user_id = $1
`

type GetReportRow struct {
	Category string
	Amount   float64
}

func (q *Queries) GetReport(ctx context.Context, userID uuid.UUID) ([]GetReportRow, error) {
	rows, err := q.db.QueryContext(ctx, getReport, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReportRow
	for rows.Next() {
		var i GetReportRow
		if err := rows.Scan(&i.Category, &i.Amount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
