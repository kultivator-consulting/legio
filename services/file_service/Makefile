.PHONY : build run

build:
	go build -o bin

run: build
	./bin

docker:
	docker build -t file_service:latest .

run-container:
	docker run --name=file_service --network="host" -d file_service

# unit test & calculate code coverage
test:
	@if [ -f test-reports/coverage.txt ]; then rm test-reports/coverage.txt; fi;
	@echo ">> running unit test and calculate coverage"
	@go test ./... -cover -coverprofile=coverage.txt -covermode=count -coverpkg=$(PACKAGES)
	@go tool cover -func=coverage.txt

clean:
	rm bin file_service
