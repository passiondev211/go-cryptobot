-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD language varchar(20);
ALTER TABLE users ADD email varchar(50);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP COLUMN language;
ALTER TABLE users DROP COLUMN email;