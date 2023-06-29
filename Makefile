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

migrateup1:
	migrate -path db/migration -database "postgresql://${USERNAME}:${PASSWORD}@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://${USERNAME}:${PASSWORD}@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://${USERNAME}:${PASSWORD}@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

create-migration:
	migrate create -ext sql -dir db/migration -seq 

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/F34th3R/go_simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock create-migration
