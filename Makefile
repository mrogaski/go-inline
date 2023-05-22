GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet

.PHONY: test

clean: ## Remove build related file
	rm -f ./cover.out

test: ## Run the tests of the project
	$(GOTEST) -v ./...

cover: ## Generate test coverage report
	$(GOTEST) -coverprofile cover.out ./...
	$(GOCMD) tool cover -html cover.out

