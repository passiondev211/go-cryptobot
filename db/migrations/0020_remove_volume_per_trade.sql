-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users DROP COLUMN volume_per_trade;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users ADD volume_per_trade DECIMAL(20,10) DEFAULT 0;
