goose -dir ./migrations postgres "postgres://user:password@localhost:5432/example?sslmode=disable" status

goose -dir ./migrations postgres "postgres://user:password@localhost:5432/example?sslmode=disable" up