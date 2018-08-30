-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user_auth (
  id int(11) NOT NULL AUTO_INCREMENT,
  user_id int(11),
  auth_token varchar(100) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  expiring_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  is_used BOOLEAN NOT NULL,
  PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE user_auth;