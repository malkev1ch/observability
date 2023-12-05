# ==============================================================================
# Main
ifneq ("$(wildcard .local/.env)","")
    include .local/.env
    export $(shell sed 's/=.*//' .local/.env)
endif

.PHONY: run-api
run-api:
	go run apiservice/main.go

run-api-root:
	ADDRESS=0.0.0.0:8080 \
	USER_SERVICE_ADDRESS=0.0.0.0:8081 \
	VOUCHER_SERVICE_ADDRESS=0.0.0.0:8082 \
	OTEL_ADDRESS=0.0.0.0:8081 \
	cd apiservice && go run main.go


COVERAGE_FILE="coverage.out"

.PHONY: test-ci
test-ci:
	go test ./... -covermode=set -coverprofile=$(COVERAGE_FILE)

.PHONY: test
test:
	go test -cover ./... -count=1

.PHONY: bench
bench:
	go test -bench=. -benchmem -count 1 -benchtime=5s ./...

.PHONY: gen
gen:
	go generate ./...

# ==============================================================================
# Deploy commands

# ==============================================================================
# Tools commands

linter:
	golangci-lint run ./...

linter-fix:
	golangci-lint run ./... --fix

openapi:
	oapi-codegen -generate types,server,spec -o apiservice/gen/v1/openapi.gen.go -package gen apiservice/api/v1/spec.yaml

buf:
	buf generate userservice/proto
	buf generate voucherservice/proto

breaking:
	buf breaking userservice/proto --against "../../.git#subdir=start/getting-started-with-buf-cli/proto"
	buf breaking voucherservice/proto --against "../../.git#subdir=start/getting-started-with-buf-cli/proto"


SWAGGER_UI_VERSION := v5.10.3
## Update assets for Swagger UI
update-swaggergui:
	curl https://raw.githubusercontent.com/swagger-api/swagger-ui/$(SWAGGER_UI_VERSION)/dist/swagger-ui-bundle.js -o ./pkg/swaggergui/static/swagger-ui-bundle.js --create-dirs
	curl https://raw.githubusercontent.com/swagger-api/swagger-ui/$(SWAGGER_UI_VERSION)/dist/swagger-ui-standalone-preset.js -o ./pkg/swaggergui/static/swagger-ui-standalone-preset.js
	curl https://raw.githubusercontent.com/swagger-api/swagger-ui/$(SWAGGER_UI_VERSION)/dist/swagger-ui.js -o ./pkg/swaggergui/static/swagger-ui.js
	curl https://raw.githubusercontent.com/swagger-api/swagger-ui/$(SWAGGER_UI_VERSION)/dist/swagger-ui.css -o ./pkg/swaggergui/static/swagger-ui.css
	curl https://raw.githubusercontent.com/swagger-api/swagger-ui/$(SWAGGER_UI_VERSION)/dist/oauth2-redirect.html -o ./pkg/swaggergui/static/oauth2-redirect.html
	curl https://raw.githubusercontent.com/swagger-api/swagger-ui/$(SWAGGER_UI_VERSION)/dist/favicon-32x32.png -o ./pkg/swaggergui/static/favicon-32x32.png
	curl https://raw.githubusercontent.com/swagger-api/swagger-ui/$(SWAGGER_UI_VERSION)/dist/favicon-16x16.png -o ./pkg/swaggergui/static/favicon-16x16.png
	rm -rf ./pkg/swaggergui/static/*.gz
	zopfli --i50 ./pkg/swaggergui/static/*.js && rm -f ./pkg/swaggergui/static/*.js
	zopfli --i50 ./pkg/swaggergui/static/*.css && rm -f ./pkg/swaggergui/static/*.css
	zopfli --i50 ./pkg/swaggergui/static/*.html && rm -f ./pkg/swaggergui/static/*.html

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy

tidy:
	go mod tidy

deps-upgrade:
  # go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d ./...
	go mod tidy

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Git support
delete-local-branches:
	git for-each-ref --format '%(refname:short)' refs/heads | grep -v 'develop\|qa\|master\|release/*' | xargs git branch -D