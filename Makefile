.PHONY: help
help:
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n" \
	} \
	/^[a-zA-Z_-]+:.*?##/ { \
		printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2 \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: install-dev-dependencies
install-dev-dependencies: ## install-dev-dependencies assumes Go is installed
	go install github.com/storj/ci/check-mod-tidy@latest
	go install github.com/storj/ci/check-copyright@latest
	go install github.com/storj/ci/check-large-files@latest
	go install github.com/storj/ci/check-imports@latest
	go install github.com/storj/ci/check-peer-constraints@latest
	go install github.com/storj/ci/check-atomic-align@latest
	go install github.com/storj/ci/check-errs@latest
	go install github.com/storj/ci/check-deferloop@latest
	go install github.com/storj/ci/check-downgrades@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	go install github.com/google/go-licenses@v1.6.0


.PHONY: lint
lint: ## Lint
	check-mod-tidy
	check-copyright
	check-large-files
	check-imports -race ./...
	check-peer-constraints -race
	check-atomic-align ./...
	check-errs ./...
	check-deferloop ./...
	staticcheck ./...
	golangci-lint run --print-resources-usage 
	check-downgrades
	go-licenses check ./...
	go vet ./...


.PHONY: deploy
deploy: ## deploy to Google Cloud Functions
	gcloud functions deploy RevokeAccess --runtime go121 --trigger-http --allow-unauthenticated
	gcloud functions deploy RestrictAccess --runtime go121 --trigger-http --allow-unauthenticated
	gcloud functions deploy OverrideEncryption --runtime go121 --trigger-http --allow-unauthenticated
	gcloud functions deploy RegisterAccess --runtime go121 --trigger-http --allow-unauthenticated
	gcloud functions deploy NewS3Customer --runtime go121 --trigger-http --allow-unauthenticated
