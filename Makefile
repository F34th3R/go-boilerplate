USERNAME = root
PASSWORD = secret

postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=${USERNAME} -e POSTGRES_PASSWORD=${PASSWORD} -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=${USERNAME} --owner=${USERNAME} simple_bank

dropdb:
	docker exec -it postgres15 dropdb --username=${USERNAME} simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://${USERNAME}:${PASSWORD}@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://${USERNAME}:${PASSWORD}@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/F34th3R/go_simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
