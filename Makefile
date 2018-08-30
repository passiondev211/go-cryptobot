.PHONY: prod
prod:
	docker-compose up -d tls

.PHONY: upd
upd:
	docker-compose up -d cryptobot

.PHONY: up
up:
	docker-compose up cryptobot

.PHONY: stop
stop:
	docker-compose stop

.PHONY: mysql
mysql:
	docker-compose exec db mysql -u cryptobot_user -pcryptobot_pass -D cryptobot_db

.PHONY: migrate
migrate:
	docker-compose run --rm cryptobot go run db/migrates.go ${type} ${n}
