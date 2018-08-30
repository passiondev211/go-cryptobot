-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD leverage int(4) NOT NULL DEFAULT 1;



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP COLUMN leverage;