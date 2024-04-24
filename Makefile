build:
	@go build -o bin/s3-tui

run: build
	@./bin/s3-tui

test:
	@go test -v ./...

clean:
	@rm -rf bin

