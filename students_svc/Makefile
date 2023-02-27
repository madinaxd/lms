postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:12
createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres students_svc
	
dropdb:
	docker exec -it postgres dropdb --username=postgres students_svc

create_migrate:
	migrate create -ext sql -dir db/migrations -seq init_schema

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/students_svc?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/students_svc?sslmode=disable" --verbose down

sqlc:
	docker run --rm -v "D:/go_projects/students_svc/db:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...
	
.PHONY: createdb dropdb create_migrate migrateup migratedown sqlc test