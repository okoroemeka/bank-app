DB_URL = postgres://root:mysecret@localhost:5433/simple_bank?sslmode=disable
postgres:
	docker run --name postgres15 --network bank-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

stpdb:
	docker stop postgres15

rmdb:
	docker rm postgres15

dbmigrationup:
	migrate -path db/migration -database "${DB_URL}" -verbose up

dbmigrationup1:
	migrate -path db/migration -database "${DB_URL}" -verbose up

dbmigrationdown:
	migrate -path db/migration -database "${DB_URL}" -verbose down

dbmigrationdown1:
	migrate -path db/migration -database "${DB_URL}" -verbose down 1
sqlc:
	sqlc generate
server:
	go run main.go
test:
	go test -v -cover ./...

mock:
	mockgen -package mockdbb -destination db/mock/store.go github.com/okoroemeka/simple_bank/db/sqlc Store

dkbuild:
	 docker build -t simplebank:latest .

dkserver:
	docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release simplebank:latest
up:
	docker compose up
down:
	docker compose down

.PHONY: dbmigrationup1 dbmigrationdown1 postgres createdb dropdb stpdb rmdb dbmigrationup dbmigrationdown sqlc test mock
