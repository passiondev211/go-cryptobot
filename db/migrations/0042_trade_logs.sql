-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE trade_logs (
  id                        INT(11)      NOT NULL AUTO_INCREMENT,
  mode                      VARCHAR(40),
  total_users               INT(11),
  successful_users          INT(11),
  zero_fix_factor_users     INT(11),
  below_min_balance_users   INT(11),
  outside_bound_users       INT(11),     
  start_at                  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  end_at                    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  db_time                   VARCHAR(40),
  total_time                VARCHAR(40),
  `created_at`              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
 
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE trade_logs;