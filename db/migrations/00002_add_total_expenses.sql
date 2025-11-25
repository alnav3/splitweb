-- +goose Up

-- Add total_expenses column to users table
ALTER TABLE users ADD COLUMN total_expenses DECIMAL(10,2) DEFAULT 0.00 NOT NULL;

-- Add index for better performance on total_expenses queries
CREATE INDEX idx_users_total_expenses ON users(total_expenses);

-- +goose Down

-- Drop index
DROP INDEX IF EXISTS idx_users_total_expenses;

-- Remove total_expenses column
ALTER TABLE users DROP COLUMN IF EXISTS total_expenses;