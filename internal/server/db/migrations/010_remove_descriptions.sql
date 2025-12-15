-- +goose Up
ALTER TABLE pods DROP COLUMN description;
ALTER TABLE projects DROP COLUMN description;

-- +goose Down
ALTER TABLE pods ADD COLUMN description TEXT;
ALTER TABLE projects ADD COLUMN description TEXT;
