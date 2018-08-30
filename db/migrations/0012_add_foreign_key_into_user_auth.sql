-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE FROM user_auth;
ALTER TABLE user_auth ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(outer_id);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE user_auth DROP FOREIGN KEY fk_user_id;