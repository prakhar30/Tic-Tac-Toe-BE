dockerpostgres:
	docker pull postgres

postgres:
	docker run --name postgresLatest -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=something_secret -d postgres

createdb:
	docker exec -it postgresLatest createdb --username=root --owner=root tictactoe

dropdb:
	docker exec -it postgresLatest dropdb tictactoe

migrateup:
	migrate -path db/migration -database "postgresql://root:something_secret@localhost:5432/tictactoe?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:something_secret@localhost:5432/tictactoe?sslmode=disable" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock: 
	mockgen --package mockdb --destination db/mock/store.go main/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans: 
	evans --host localhost --port 9091 -r repl

.PHONY: dockerpostgres postgres createdb dropdb migrateup migratedown sqlc test server mock proto evans
