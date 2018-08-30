-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE deposit_logs (
  id int(11) NOT NULL AUTO_INCREMENT,
  amount varchar(20) NOT NULL,
  user_id int(11) NOT NULL,
  is_first_time boolean DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE deposit_logs;