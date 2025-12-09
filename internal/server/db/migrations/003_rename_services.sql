-- +goose Up
ALTER TABLE services RENAME TO pods;

-- +goose Down
ALTER TABLE pods RENAME TO services;
