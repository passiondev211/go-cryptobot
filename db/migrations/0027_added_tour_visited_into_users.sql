-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD tour_visited boolean DEFAULT 0 NOT NULL;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP COLUMN tour_visited;
