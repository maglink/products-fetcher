build:
	go build -o ./bin/products-fetcher ./cmd/products-fetcher

run: build
	./bin/products-fetcher

test:
	go test ./...

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pkg/messages/messages.proto

example-csv-gen:
	cd tests/example_generator && go run main.go

docker-run:
	docker-compose build --no-cache
	docker-compose up --scale app=2