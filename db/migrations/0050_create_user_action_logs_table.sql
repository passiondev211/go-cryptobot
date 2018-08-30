-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user_action_logs (
  id int(11) NOT NULL AUTO_INCREMENT,
  action varchar(20) NOT NULL,
  user_id int(11) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE user_action_logs;