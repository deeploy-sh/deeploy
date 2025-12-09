ALTER TABLE pods ADD COLUMN repo_url TEXT;
ALTER TABLE pods ADD COLUMN branch TEXT DEFAULT 'main';
ALTER TABLE pods ADD COLUMN dockerfile_path TEXT DEFAULT 'Dockerfile';
ALTER TABLE pods ADD COLUMN container_id TEXT;
ALTER TABLE pods ADD COLUMN status TEXT DEFAULT 'stopped';
