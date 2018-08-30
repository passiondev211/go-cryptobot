-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users (
  id int(11) NOT NULL AUTO_INCREMENT,
  balance DECIMAL(20,10) NOT NULL,
  volume_per_trade DECIMAL(20,10) NOT NULL,
  PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE users;