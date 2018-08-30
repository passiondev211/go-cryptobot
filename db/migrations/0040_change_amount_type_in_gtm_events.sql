-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE gtm_events MODIFY amount DECIMAL(20,10) NOT NULL;
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE gtm_events MODIFY amount INT(1) NOT NULL;