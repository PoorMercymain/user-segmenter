BEGIN TRANSACTION;
CREATE TABLE deletion_times (user_id TEXT, slug TEXT, deletion_timestamp TIMESTAMP WITH TIME ZONE, PRIMARY KEY(user_id, slug));
CREATE INDEX deletion_times_idx ON deletion_times USING BTREE (user_id, slug, deletion_timestamp);
CREATE TABLE users_segment_history (user_id TEXT, slug TEXT, modified_at TIMESTAMP WITH TIME ZONE, is_deletion BOOLEAN);
CREATE TABLE slugs (slug TEXT PRIMARY KEY);
CREATE INDEX slugs_idx ON slugs USING BTREE (slug);
CREATE TABLE users (user_id TEXT PRIMARY KEY, slugs TEXT[]);
CREATE INDEX users_idx ON users USING BTREE (user_id);
COMMIT;