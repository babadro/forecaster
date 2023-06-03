package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lann/builder"
)
import sq "github.com/Masterminds/squirrel"

type ForecastDB struct {
	db *pgxpool.Pool
	q  sq.StatementBuilderType
}

func NewForecastDB(db *pgxpool.Pool) *ForecastDB {
	return &ForecastDB{
		db: db,
		q:  sq.StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(sq.Dollar),
	}
}
