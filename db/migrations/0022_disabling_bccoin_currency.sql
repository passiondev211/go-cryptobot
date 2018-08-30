-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM currencies WHERE name = 'BCCOIN';



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
INSERT INTO currencies (name) VALUES ('BCCOIN');
