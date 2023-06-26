include .env

migrate:
	migrate -database="${TE_DB_DSN}" -path=migrations/ up

drop:
	migrate -database="${TE_DB_DSN}" -path=migrations/ drop

reset: drop migrate

