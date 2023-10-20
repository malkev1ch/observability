# ==============================================================================
# Main
ifneq ("$(wildcard .local/.env)","")
    include .local/.env
    export $(shell sed 's/=.*//' .local/.env)
endif

.PHONY: run-api
run-api:
	go run apiservice/main.go

.PHONY: build
build:
	go build main.go

.PHONY: test
test:
	go test -cover ./... -count=1

.PHONY: gen
gen:
	go generate ./...

# ==============================================================================
# Deploy commands

.PHONY: build-and-push-image
build-and-push-image:
	docker build -t europe-docker.pkg.dev/tactical-works-402510/registry/apiservice:latest -f apiservice/Dockerfile
	docker push europe-docker.pkg.dev/tactical-works-402510/registry/apiservice

.PHONY: deploy
deploy:
	helm upgrade otel-collector-ds open-telemetry/opentelemetry-collector --values deploy/optl/daemonset.yaml
	#helm upgrade otel-collector-dp open-telemetry/opentelemetry-collector --values deploy/optl/deployment.yaml

# ==============================================================================
# Tools commands

linter:
	golangci-lint run ./...

linter-fix:
	golangci-lint run ./... --fix

buf:
	buf generate userservice/proto

breaking:
	buf breaking userservice/proto --against "../../.git#subdir=start/getting-started-with-buf-cli/proto"

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