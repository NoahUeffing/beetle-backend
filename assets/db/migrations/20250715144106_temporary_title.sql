-- +goose Up
-- Drop the existing unique constraint on username
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;

-- Create a partial unique index to enforce uniqueness only for non-deleted users
CREATE UNIQUE INDEX unique_active_username ON users (username) WHERE deleted_at IS NULL;

-- +goose Down
-- Drop the partial unique index
DROP INDEX IF EXISTS unique_active_username;

-- Step 1: Remove or resolve duplicate usernames among soft-deleted users
WITH duplicates AS (
    SELECT username
    FROM users 
    GROUP BY username 
    HAVING COUNT(*) > 1
),
conflicts AS (
    SELECT u.id, u.username,
           ROW_NUMBER() OVER (PARTITION BY u.username ORDER BY u.created_at) AS rn
    FROM users u
    JOIN duplicates d ON u.username = d.username
)
UPDATE users
SET username = users.username || '_' || users.id::text
FROM conflicts
WHERE users.id = conflicts.id
  AND conflicts.rn > 1
  AND users.deleted_at IS NOT NULL;

-- Step 2: Restore original unique constraint on username
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);