package db

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	models "github.com/babadro/forecaster/internal/models/swagger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lann/builder"
)

type TestDB struct {
	DB *pgxpool.Pool
	q  sq.StatementBuilderType
}

func NewTestDB(db *pgxpool.Pool) *TestDB {
	return &TestDB{
		DB: db,
		q:  sq.StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(sq.Dollar),
	}
}

func (db *TestDB) CreateSeries(ctx context.Context, s models.Series) (res models.Series, err error) {
	now := time.Now()

	seriesSQL, args, err := db.q.
		Insert("forecaster.series").Columns("id", "title", "description", "updated_at", "created_at").
		Values(s.ID, s.Title, s.Description, now, now).
		Suffix("RETURNING id, title, description, created_at, updated_at").
		ToSql()

	if err != nil {
		return models.Series{}, buildingQueryFailed("insert series", err)
	}

	err = db.DB.QueryRow(ctx, seriesSQL, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.UpdatedAt, &res.CreatedAt)
	if err != nil {
		return models.Series{}, scanFailed("insert series", err)
	}

	return res, nil
}

func buildingQueryFailed(queryName string, err error) error {
	return fmt.Errorf("%s: building query failed: %w", queryName, err)
}

func queryFailed(queryName string, err error) error {
	return fmt.Errorf("%s: query failed: %w", queryName, err)
}

func rowsError(queryName string, err error) error {
	return fmt.Errorf("%s: rows error: %w", queryName, err)
}

func scanFailed(queryName string, err error) error {
	return fmt.Errorf("%s: scan rows failed: %w", queryName, err)
}

func execFailed(queryName string, err error) error {
	return fmt.Errorf("%s: exec failed: %w", queryName, err)
}
