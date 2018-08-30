-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM currencies WHERE name = 'BCH';
DELETE FROM currencies WHERE name = 'BCCOIN';
DELETE FROM currencies WHERE name = 'XMR';

INSERT INTO currencies (name) VALUES ('XRP');
INSERT INTO currencies (name) VALUES ('BCC');
INSERT INTO currencies (name) VALUES ('ZEC');
INSERT INTO currencies (name) VALUES ('ETC');


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
INSERT INTO currencies (name) VALUES ('BCH');;
INSERT INTO currencies (name) VALUES ('BCCOIN');
INSERT INTO currencies (name) VALUES ('XMR');

DELETE FROM currencies WHERE name = 'XRP';
DELETE FROM currencies WHERE name = 'BCC';
DELETE FROM currencies WHERE name = 'ZEC';
DELETE FROM currencies WHERE name = 'ETC';


