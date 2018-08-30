-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD daily_trade_profit_bounds_lower DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD daily_trade_profit_bounds_upper DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD result_trade_profit_bounds_lower DECIMAL(20,10) NOT NULL DEFAULT 0;
ALTER TABLE users ADD result_trade_profit_bounds_upper DECIMAL(20,10) NOT NULL DEFAULT 0;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users DROP COLUMN daily_trade_profit_bounds_lower;
ALTER TABLE users DROP COLUMN daily_trade_profit_bounds_upper;
ALTER TABLE users DROP COLUMN result_trade_profit_bounds_lower;
ALTER TABLE users DROP COLUMN result_trade_profit_bounds_upper;
