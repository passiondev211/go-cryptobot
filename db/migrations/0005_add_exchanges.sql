-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO exchanges (name) VALUES ('Bitfinex');
INSERT INTO exchanges (name) VALUES ('Bithumb');
INSERT INTO exchanges (name) VALUES ('coinone');
INSERT INTO exchanges (name) VALUES ('Bittrex');
INSERT INTO exchanges (name) VALUES ('HitBTC');
INSERT INTO exchanges (name) VALUES ('Poloniex');
INSERT INTO exchanges (name) VALUES ('Quoine');
INSERT INTO exchanges (name) VALUES ('Coinbase GDAX');
INSERT INTO exchanges (name) VALUES ('Kraken');
INSERT INTO exchanges (name) VALUES ('Bitstamp');
INSERT INTO exchanges (name) VALUES ('bitFlyer');
INSERT INTO exchanges (name) VALUES ('Binance');
INSERT INTO exchanges (name) VALUES ('Gemini');
INSERT INTO exchanges (name) VALUES ('Korbit');


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM exchanges WHERE name = 'Bitfinex';
DELETE FROM exchanges WHERE name = 'Bithumb';
DELETE FROM exchanges WHERE name = 'coinone';
DELETE FROM exchanges WHERE name = 'Bittrex';
DELETE FROM exchanges WHERE name = 'HitBTC';
DELETE FROM exchanges WHERE name = 'Poloniex';
DELETE FROM exchanges WHERE name = 'Quoine';
DELETE FROM exchanges WHERE name = 'Coinbase GDAX';
DELETE FROM exchanges WHERE name = 'Kraken';
DELETE FROM exchanges WHERE name = 'Bitstamp';
DELETE FROM exchanges WHERE name = 'bitFlyer';
DELETE FROM exchanges WHERE name = 'Binance';
DELETE FROM exchanges WHERE name = 'Gemini';
DELETE FROM exchanges WHERE name = 'Korbit';


