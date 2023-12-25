DB_URL = postgres://root:mysecret@postgres:5432/simple_bank?sslmode=disable

postgres:
	docker network rm bank-network  &&	docker network create bank-network && docker run --name postgres15 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

stpdb:
	docker stop postgres15

rmdb:
	docker rm postgres15 && docker network rm bank-network

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

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto && \
    statik -src=./doc/swagger -dest=./doc -f

evans:
	evans --host localhost --port 9080 -r repl

.PHONY: dbmigrationup1 dbmigrationdown1 postgres createdb dropdb stpdb rmdb dbmigrationup dbmigrationdown sqlc test mock dkbuild dkserver up down proto
