-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM users WHERE id = 1;
ALTER TABLE user_auth ADD CONSTRAINT auth_token_uniq UNIQUE (auth_token);
ALTER TABLE trades CHANGE buy_exchange buy_exchange VARCHAR(128) NOT NULL;
ALTER TABLE trades CHANGE sell_exchange sell_exchange VARCHAR(128) NOT NULL;
ALTER TABLE exchanges CHANGE name name VARCHAR(128) NOT NULL;
ALTER TABLE user_auth CHANGE auth_token auth_token VARCHAR(256) NOT NULL;
ALTER TABLE users CHANGE email email VARCHAR(128) NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
INSERT INTO users (balance, outer_id, email) VALUES (2.0, 0, 'test');
ALTER TABLE user_auth DROP INDEX `auth_token_uniq`;
ALTER TABLE trades CHANGE buy_exchange buy_exchange VARCHAR(20) NOT NULL;
ALTER TABLE trades CHANGE sell_exchange sell_exchange VARCHAR(20) NOT NULL;
ALTER TABLE exchanges CHANGE name name VARCHAR(20) NOT NULL;
ALTER TABLE user_auth CHANGE auth_token auth_token VARCHAR(100) NOT NULL;
ALTER TABLE users CHANGE email email VARCHAR(50) NOT NULL;