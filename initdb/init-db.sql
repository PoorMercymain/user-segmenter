BEGIN TRANSACTION;
CREATE TABLE deletion_times (user_id TEXT, slug TEXT, deletion_timestamp TIMESTAMP WITH TIME ZONE);
CREATE TABLE users_segment_history (user_id TEXT, slug TEXT, modified_at TIMESTAMP WITH TIME ZONE, is_deletion BOOLEAN);
CREATE TABLE slugs (slug TEXT PRIMARY KEY);
CREATE TABLE users (user_id TEXT PRIMARY KEY, slugs TEXT[]);
COMMIT;