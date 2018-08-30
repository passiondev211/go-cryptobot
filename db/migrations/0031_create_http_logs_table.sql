-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE http_logs (
  id int(11) NOT NULL AUTO_INCREMENT,
  client_ip varchar(40) NOT NULL,
  url text NOT NULL,
  request text,
  response text,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE http_logs;