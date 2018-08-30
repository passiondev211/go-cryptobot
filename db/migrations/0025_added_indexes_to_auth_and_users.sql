-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX outer_id_index ON users(outer_id);
CREATE INDEX auth_token_index ON user_auth(auth_token);
CREATE INDEX is_used_index ON user_auth(is_used);
CREATE INDEX expiring_at_index ON user_auth(expiring_at);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX outer_id_index ON users;
DROP INDEX auth_token_index ON user_auth;
DROP INDEX is_used_index ON user_auth;
DROP INDEX expiring_at_index ON user_auth;
