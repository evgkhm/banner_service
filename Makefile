.PHONY: \
	run \
	lint \
	mocks \
	test \

run:
	docker-compose up --build

lint:
	golangci-lint cache clean
	golangci-lint run --config=./.golangci.yaml

test: ### run test
	go test -v ./...

MOCKS_DESTINATION=mocks
mocks: internal/usecase/deps.go internal/usecase/usecase.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done