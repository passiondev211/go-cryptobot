-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM exchanges WHERE name = 'Bithumb';
DELETE FROM exchanges WHERE name = 'coinone';
DELETE FROM exchanges WHERE name = 'Quoine';
DELETE FROM exchanges WHERE name = 'Coinbase GDAX';
DELETE FROM exchanges WHERE name = 'Bitstamp';
DELETE FROM exchanges WHERE name = 'bitFlyer';
DELETE FROM exchanges WHERE name = 'Gemini';
DELETE FROM exchanges WHERE name = 'Korbit';

INSERT INTO exchanges (name) VALUES ('Huobi');
INSERT INTO exchanges (name) VALUES ('OKEx');
INSERT INTO exchanges (name) VALUES ('Upbit');


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
INSERT INTO exchanges (name) VALUES ('Bithumb');
INSERT INTO exchanges (name) VALUES ('coinone');;
INSERT INTO exchanges (name) VALUES ('Quoine');
INSERT INTO exchanges (name) VALUES ('Coinbase GDAX');
INSERT INTO exchanges (name) VALUES ('Bitstamp');
INSERT INTO exchanges (name) VALUES ('bitFlyer');
INSERT INTO exchanges (name) VALUES ('Gemini');
INSERT INTO exchanges (name) VALUES ('Korbit');

DELETE FROM exchanges WHERE name = 'Huobi';
DELETE FROM exchanges WHERE name = 'OKEx';
DELETE FROM exchanges WHERE name = 'Upbit';


