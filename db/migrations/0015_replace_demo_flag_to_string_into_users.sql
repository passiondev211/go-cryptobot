-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD mode varchar(20) DEFAULT 'demo';
UPDATE users SET mode='real' WHERE demo=0;
ALTER TABLE users DROP COLUMN demo;



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE users ADD demo BOOLEAN NOT NULL DEFAULT 1;
UPDATE users SET demo=0 WHERE mode='real';
ALTER TABLE users DROP COLUMN mode;
