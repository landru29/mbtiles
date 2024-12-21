build:
	go build -o mbtiles ./cmd

clean:
	rm -rf mbtiles

lint:
	golangci-lint run ./...

test:
	gotest ./...

generate:
	go generate ./...