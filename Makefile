build:
	docker build -t backend:1 --file ./api/Dockerfile .
	docker-compose build
run: build
	docker-compose up -d
stop:
	docker-compose down