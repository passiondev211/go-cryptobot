-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE http_logs ADD response_time VARCHAR(255);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE http_logs DROP COLUMN response_time;
