-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE trades ADD profit DECIMAL(20,10) NOT NULL;



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE trades DROP COLUMN profit;