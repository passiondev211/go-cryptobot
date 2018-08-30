-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE notifications (
  id           INT(11)      NOT NULL AUTO_INCREMENT,
  `user`       INT(11)      NOT NULL,
  `status`     INT(1)       NOT NULL,
  `message`    VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE notifications;