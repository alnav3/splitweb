package databaselogic

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	internalRepository "github.com/alnav3/splitweb/db/internal_repository"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func SetupDatabase() (*Repository, error) {
    ctx := context.Background()
	url := os.Getenv("DATABASE_URL")

	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	// Create pgx pool for queries
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Create sql.DB for goose migrations
	db, err := sql.Open("pgx", url)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to open sql.DB for migrations: %w", err)
	}

	// Run migrations with goose
	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("db/migrations"))
	if err != nil {
		pool.Close()
		db.Close()
		return nil, fmt.Errorf("error creating goose provider: %w", err)
	}

	results, err := provider.Up(context.Background())
	if err != nil {
		pool.Close()
		db.Close()
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	for _, r := range results {
		log.Printf("OK   %s (%s)\n", r.Source.Path, r.Duration)
	}

	// Create queries with the pgx pool
    queries := internalRepository.New(pool)

    return &Repository{
        Queries: queries,
        Pool:    pool,
        DB:      db,
        Context: ctx,
    }, nil
}

// PerformDBOperations is an example function for database operations
func PerformDBOperations(ctx context.Context, repo *Repository) error {
    // Example: get all users
    users, err := repo.Queries.FindAllUsers(ctx)
    if err != nil {
        return fmt.Errorf("error finding users: %w", err)
    }

    log.Printf("Found %d users", len(users))
    for _, user := range users {
        log.Printf("User: %s (%s)", user.Email, user.ID)
    }

    return nil
}
