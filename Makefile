.PHONY: git-init git-checkout golangci-lint-install lint test yc-zip yc-create-function-or-ignore yc-create-function-version yc-timer yc-clear yc-deploy

ifneq (,$(wildcard ./.env))
include .env
export
endif

REPO_NAME := $(shell basename $(CURDIR))
PROJECT := $(CURDIR)
LOCAL_BIN := $(CURDIR)/bin

# GIT
git-init:
	echo '/.idea/\n/bin/\n.env\n' > .gitignore
	gh repo create $(GIT_USER)/$(REPO_NAME) --private
	git init
	git config user.name "$(GIT_USER)"
	git config user.email "$(GIT_EMAIL)"
	git add Makefile go.mod .gitignore
	git commit -m "Init commit"
	git remote add origin git@github.com:$(GIT_USER)/$(REPO_NAME).git
	git remote -v
	git push -u origin master

BN ?= dev
git-checkout:
	git checkout -b $(BN)

# LINT
lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

lint:
	$(LOCAL_BIN)/golangci-lint run ./...
	

# TEST
test:
	go vet ./...
	go test -v ./internal/...

# Test server endpoints
test-success:
	# Test the successful endpoint, expecting a 200 OK response
	curl -i -X PUT -H "Content-Type: application/json" -d '{"api_key": "KEY", "timestamp":1744761600}' http://localhost:8080/api/updatePriceNew | grep "HTTP/1.1 200 OK"

test-error:
	# Test the error endpoint, expecting a 500 Internal Server Error response
	curl -i -X PUT http://localhost:8080/api/updatePriceNew/error500 | grep "HTTP/1.1 500 Internal Server Error"

test-noresponse:
	# Test the no response endpoint. Using timeout to handle no response.
	curl -m 2 -X PUT http://localhost:8080/api/updatePriceNew/noresponse || true

# YANDEX CLOUD
yc-zip:
		zip -r '$(YCF_FUNC_NAME).zip' handler.go go.mod go.sum internal

yc-create-function-or-ignore:
	@if yc serverless function get --name '$(YCF_FUNC_NAME)' > /dev/null 2>&1; then \
		echo "Function '$(YCF_FUNC_NAME)' already exists, skipping creation."; \
	else \
		yc serverless function create --name "$(YCF_FUNC_NAME)"; \
	fi

yc-create-function-version:
	yc serverless function version create \
	--function-name '$(YCF_FUNC_NAME)' \
	--service-account-id '$(YCF_SA_ID)' \
	--runtime golang121 \
	--entrypoint handler.Handler \
	--execution-timeout $(APP_TIMEOUT_DURATION)s \
	--memory 128m \
	--environment APP_MODE=prod \
	--source-path "./$(YCF_FUNC_NAME).zip"
	
yc-timer:
	@if yc serverless trigger get --name "run-$(YCF_FUNC_NAME)" > /dev/null 2>&1; then \
  		echo "Trigger 'run-$(YCF_FUNC_NAME)' already exists, skipping creation."; \
	else \
		yc serverless trigger create timer \
		--cron-expression $(YCF_CRON) \
		--invoke-function-name "$(YCF_FUNC_NAME)" \
		--invoke-function-service-account-id "$(YCF_SA_ID)" \
		--name "run-$(YCF_FUNC_NAME)"; \
	fi

yc-clear:
	rm "./$(YCF_FUNC_NAME).zip"

yc-deploy: test yc-zip yc-create-function-or-ignore yc-create-function-version yc-timer yc-clear
