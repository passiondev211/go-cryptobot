-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD country VARCHAR(3);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP COLUMN country;
