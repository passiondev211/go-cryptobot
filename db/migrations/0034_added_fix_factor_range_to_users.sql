-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users DROP COLUMN fix_factor_min_percentage;
ALTER TABLE users DROP COLUMN fix_factor_max_percentage;
ALTER TABLE users CHANGE fix_factor_min_value min_fix_factor DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users CHANGE fix_factor_max_value max_fix_factor DECIMAL(20,10) NOT NULL DEFAULT 0;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users ADD fix_factor_min_percentage DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD fix_factor_max_percentage DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users CHANGE min_fix_factor fix_factor_min_value DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users CHANGE max_fix_factor fix_factor_max_value DECIMAL(20,10) NOT NULL DEFAULT 0;
