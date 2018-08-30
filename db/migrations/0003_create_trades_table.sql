-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE trades (
  id int(11) NOT NULL AUTO_INCREMENT,
  currency varchar(20) NOT NULL,
  buy_price DECIMAL(20,10) NOT NULL,
  sell_price DECIMAL(20,10) NOT NULL,
  buy_exchange varchar(20) NOT NULL,
  sell_exchange varchar(20) NOT NULL,
  amount DECIMAL(20,10) NOT NULL,
  created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE trades;