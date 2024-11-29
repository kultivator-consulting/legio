.PHONY : prepare build run

$(eval $(service):;@:)

check:
	@[ "${service}" ] || ( echo "\x1b[31;1mERROR: 'service' is not set\x1b[0m"; exit 1 )
	@if [ ! -d "services/$(service)" ]; then  echo "\x1b[31;1mERROR: service '$(service)' undefined\x1b[0m"; exit 1; fi

prepare: check
	@if [ ! -f services/$(service)/.env ]; then cp services/$(service)/.env.sample services/$(service)/.env; fi;

build: check
	@go build -o services/$(service)/bin services/$(service)/*.go

run: build
	@chmod +x services/$(service)/bin
	@cd services/$(service) && ./bin

build-docker: check
	@echo "\x1b[32;1m>>> building docker image for service $(service)\x1b[0m"
	docker build -f services/$(service)/service.dockerfile --build-arg SERVICE_NAME=$(service) -t $(service):latest .

build-all-docker:
	@echo "\x1b[32;1m>>> building docker image for all services\x1b[0m"
	@for s in $(shell ls ./services); do \
		echo "\x1b[32;1m>>> building docker image for service $${s}\x1b[0m"; \
		docker build -f services/$${s}/service.dockerfile --build-arg SERVICE_NAME=$${s} -t $${s}:latest .; \
  	done

export-docker-images:
	@echo "\x1b[32;1m>>> exporting docker images for all services\x1b[0m"
	@for s in $(shell ls ./services); do \
		echo "\x1b[32;1m>>> exporting docker image for service $${s}\x1b[0m"; \
		docker save -o ./export_images/$${s}.tar $${s}:latest; \
  	done

import-docker-images:
	@echo "\x1b[32;1m>>> importing docker images for all services\x1b[0m"
	@for s in $(shell ls ./services); do \
		echo "\x1b[32;1m>>> importing docker image for service $${s}\x1b[0m"; \
		docker load -i ./import-images/$${s}.tar; \
  	done

clean-imported-images:
	@echo "\x1b[32;1m>>> removing docker images for all services\x1b[0m"
	@for s in $(shell ls ./services); do \
		echo "\x1b[32;1m>>> removing docker image for service $${s}\x1b[0m"; \
		docker image rm $${s}:latest; \
  	done

run-docker:
	docker run --name=$(service) --network="host" -d $(service):latest

test: check
	@echo "\x1b[32;1m>>> running unit test and calculate coverage for service $(service)\x1b[0m"
	@if [ ! -d ./test-reports ]; then mkdir -p ./test-reports; fi
	@if [ -f ./test-reports/$(service)_coverage.txt ]; then rm ./test-reports/$(service)_coverage.txt; fi
	@go test -race ./services/$(service)/... -cover -coverprofile=./test-reports/$(service)_coverage.txt -covermode=atomic \
		-coverpkg=$$(go list ./services/$(service)/... | grep -v -e mocks -e codebase | tr '\n' ',')
	@go tool cover -func=./test-reports/$(service)_coverage.txt

clean: check
	rm -rf services/$(service)/bin bin
