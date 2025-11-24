package databaselogic

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	internalRepository "github.com/alnav3/splitweb/db/internal_repository"
)

type Repository struct {
    Queries *internalRepository.Queries
	Pool    *pgxpool.Pool
	DB      *sql.DB  // Required for goose migrations
    Context context.Context
}

func (r *Repository) Close() {
	if r.Pool != nil {
		r.Pool.Close()
	}
	if r.DB != nil {
		r.DB.Close()
	}
}

