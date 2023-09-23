postgres:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

stpdb:
	docker stop postgres15

rmdb:
	docker rm postgres15

dbmigrationup:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5433/simple_bank?sslmode=disable" -verbose up

dbmigrationdown:
	migrate -path db/migration -database "postgresql://root:mysecret@localhost:5433/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...


.PHONY: postgres createdb dropdb stpdb rmdb dbmigrationup dbmigrationdown sqlc test
