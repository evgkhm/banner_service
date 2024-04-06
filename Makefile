.PHONY: \
	run \
	lint \

run:
	docker-compose up --build

lint:
	golangci-lint cache clean
	golangci-lint run --config=./.golangci.yaml