-- +goose Up
ALTER TABLE pods ADD COLUMN repo_url TEXT;
ALTER TABLE pods ADD COLUMN branch TEXT DEFAULT 'main';
ALTER TABLE pods ADD COLUMN dockerfile_path TEXT DEFAULT 'Dockerfile';
ALTER TABLE pods ADD COLUMN container_id TEXT;
ALTER TABLE pods ADD COLUMN status TEXT DEFAULT 'stopped';

-- +goose Down
ALTER TABLE pods DROP COLUMN repo_url;
ALTER TABLE pods DROP COLUMN branch;
ALTER TABLE pods DROP COLUMN dockerfile_path;
ALTER TABLE pods DROP COLUMN container_id;
ALTER TABLE pods DROP COLUMN status;
