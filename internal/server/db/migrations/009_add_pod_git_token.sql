-- +goose Up
ALTER TABLE pods ADD COLUMN git_token_id TEXT REFERENCES git_tokens(id);

-- +goose Down
ALTER TABLE pods DROP COLUMN git_token_id;
