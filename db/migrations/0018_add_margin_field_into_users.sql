-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD max_margin DECIMAL(20,10) DEFAULT 0.25;



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE users DROP COLUMN max_margin;