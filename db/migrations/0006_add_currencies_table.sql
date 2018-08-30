-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE currencies (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(20) NOT NULL,
  PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE currencies;