-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX user_index ON trades(user);
CREATE INDEX created_at_index ON trades(created_at);
CREATE INDEX mode_index ON trades(mode);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX user_index ON trades;
DROP INDEX created_at_index ON trades;
DROP INDEX mode_index ON trades;
