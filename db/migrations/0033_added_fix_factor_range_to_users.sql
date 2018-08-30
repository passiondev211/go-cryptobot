-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD fix_factor_min_percentage DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD fix_factor_max_percentage DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD fix_factor_min_value DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD fix_factor_max_value DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD has_custom_fix_factor boolean NOT NULL DEFAULT 0;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP COLUMN fix_factor_min_percentage;
ALTER TABLE users DROP COLUMN fix_factor_max_percentage;
ALTER TABLE users DROP COLUMN fix_factor_min_value;
ALTER TABLE users DROP COLUMN fix_factor_max_value;
ALTER TABLE users DROP COLUMN has_custom_fix_factor;
