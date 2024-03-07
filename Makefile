include .env

run:
	go run ./cmd/blogagg

gup:
	cd sql/schema && goose postgres ${DBURL} up

gdown:
	cd sql/schema && goose postgres ${DBURL} down

gstatus:
	cd sql/schema && goose postgres ${DBURL} status

sqlc:
	sqlc generate

lint:
	golangci-lint run
