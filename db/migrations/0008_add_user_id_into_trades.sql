-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE trades ADD user int(11) NOT NULL REFERENCES users(id);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE trades DROP COLUMN user;