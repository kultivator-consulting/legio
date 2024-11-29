.PHONY : build run

build:
	go build -o bin

run: build
	./bin

docker:
	docker build -t seed_command:latest .

run-container:
	docker run --name=seed_command --network="host" -d seed_command

clean:
	rm bin seed_command
