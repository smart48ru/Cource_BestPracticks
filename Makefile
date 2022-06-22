.PHONY: build
build:
	@echo 'Start build'
	@go build -o bin/tgbot cmd/smart48bot/main.go
	@echo 'The app was successfully built at bin/tgbot'

.PHONY: run
run: build
	@echo 'Start app ./bin/tgbot'
	@bin/tgbot

.PHONY: test_cover
test_cover:
	@echo 'Start Unit test cover'
	@go test ./... -coverprofile fmt

.PHONY: test
test:
	@echo 'Start Unit test'
	@go test -race -covermode=atomic ./...

.PHONY: integration
integration:
	@echo 'Start Integration test'
	@go test -race -covermode=atomic --tags=integration ./...

.PHONY: lint
lint:
	@echo 'Start lint'
	golangci-lint run

MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: internal/telegram_bot/telegram_bot.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
