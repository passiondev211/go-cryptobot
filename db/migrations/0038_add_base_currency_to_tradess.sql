-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE trades ADD base_currency VARCHAR(5) DEFAULT 'BTC' NOT NULL;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE trades DROP COLUMN base_currency;
