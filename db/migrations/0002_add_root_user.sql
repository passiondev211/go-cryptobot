-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO users (balance, volume_per_trade) VALUES (2.0, 0.01);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM users WHERE id = 1;