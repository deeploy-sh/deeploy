-- +goose Up
ALTER TABLE pod_domains ADD COLUMN type TEXT DEFAULT 'custom';
ALTER TABLE pod_domains ADD COLUMN port INTEGER DEFAULT 80;

-- +goose Down
ALTER TABLE pod_domains DROP COLUMN type;
ALTER TABLE pod_domains DROP COLUMN port;
