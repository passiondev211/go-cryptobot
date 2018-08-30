-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE trades ADD fix_factor DECIMAL(20,10) NOT NULL DEFAULT 0;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE trades DROP COLUMN fix_factor;
