-- Drop indexes
DROP INDEX IF EXISTS post_id_idx;
DROP INDEX IF EXISTS user_id_idx;

-- Drop foreign key constraints
ALTER TABLE comments DROP CONSTRAINT IF EXISTS user_id_fk;
ALTER TABLE comments DROP CONSTRAINT IF EXISTS post_id_fk;

-- Drop the table
DROP TABLE IF EXISTS comments;
