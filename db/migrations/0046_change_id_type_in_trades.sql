-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE trades MODIFY id BIGINT(20) NOT NULL AUTO_INCREMENT;
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE trades MODIFY id INT(11) NOT NULL;