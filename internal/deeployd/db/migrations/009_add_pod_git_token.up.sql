ALTER TABLE pods ADD COLUMN git_token_id TEXT REFERENCES git_tokens(id);
