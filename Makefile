include .env

gup:
	cd sql/schema && goose postgres ${DBURL} up

gdown:
	cd sql/schema && goose postgres ${DBURL} down

gstatus:
	cd sql/schema && goose postgres ${DBURL} status

lint:
	golangci-lint run
