-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO currencies (name) VALUES ('BCH');
INSERT INTO currencies (name) VALUES ('ETH');
INSERT INTO currencies (name) VALUES ('DASH');
INSERT INTO currencies (name) VALUES ('BCCOIN');
INSERT INTO currencies (name) VALUES ('XMR');
INSERT INTO currencies (name) VALUES ('LTC');



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM currencies WHERE name = 'BCH';
DELETE FROM currencies WHERE name = 'ETH';
DELETE FROM currencies WHERE name = 'DASH';
DELETE FROM currencies WHERE name = 'BCCOIN';
DELETE FROM currencies WHERE name = 'XMR';
DELETE FROM currencies WHERE name = 'LTC';


