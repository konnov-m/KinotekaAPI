build:
	docker build -t backend:1 --file ./api/Dockerfile .
	docker build -t migrate:1 --file ./migrate/Dockerfile .
	docker-compose build
run: build
	docker-compose up -d
stop:
	docker-compose down
generate:
	cd ./api/internal/service && go generate
test: generate
	cd ./api && go test ./... -v