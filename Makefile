DOCKER_IMAGE ?= $(notdir $(CURDIR))-toolchain:local
DOCKER_GATE ?= ./scripts/docker/docker_gate.sh
GO_VERSION ?= 1.24
GOLANGCI_LINT_VERSION ?= v2.1.6
GOVULNCHECK_VERSION ?= v1.1.4

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -coverpkg=./pkg/... $$(go list ./... | grep -v internal/tools)

.PHONY: race
race:
	go test -race -coverpkg=./pkg/... $$(go list ./... | grep -v internal/tools)

.PHONY: build
build:
	go build ./...

.PHONY: build-check
build-check: build

.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed"; \
		exit 1; \
	fi

.PHONY: integration
integration:
	./scripts/run_integration.sh

.PHONY: security
security:
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not installed"; \
		exit 1; \
	fi
	./scripts/check_secrets.sh

.PHONY: boundary
boundary:
	./scripts/check_boundary.sh

.PHONY: contracts
contracts:
	./scripts/check_contracts.sh

.PHONY: runtime-check
runtime-check: test vet

.PHONY: drift-check
drift-check:
	./scripts/docker/check_contract.sh

.PHONY: goalcli
goalcli:
	@echo "testkitx has no goalcli binary; L1 compatibility target only"

.PHONY: goalcli-image
goalcli-image: build-check

.PHONY: goalcli-version
goalcli-version:
	@echo "testkitx-no-goalcli"

.PHONY: property
property:
	go test ./... -run 'Test.*Property|Test.*Invariant'

.PHONY: fuzz-smoke
fuzz-smoke:
	./scripts/run_fuzz_smoke.sh

.PHONY: golden
golden:
	go test ./... -run 'Test.*Golden|Test.*Snapshot'

.PHONY: manifest-fixture-check
manifest-fixture-check:
	go test ./pkg/testkitx/manifesttest/... -run 'Test.*Manifest|Test.*Checksum'

.PHONY: downstream-test-only-adoption-check
downstream-test-only-adoption-check:
	./scripts/check_downstream_adoption.sh

.PHONY: evidence
evidence:
	./scripts/generate_manifest.sh

.PHONY: release-evidence-check
release-evidence-check:
	RELEASE_EVIDENCE_REQUIRE_PASSED=1 ./scripts/check_release_evidence.sh

.PHONY: coverage-check
coverage-check:
	@go test -coverpkg=./pkg/... $$(go list ./... | grep -v internal/tools) -coverprofile=coverage.out -covermode=atomic -timeout 180s
	@COVERAGE=$$(go tool cover -func=coverage.out | grep 'total:' | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$COVERAGE%"; \
	if [ $$(echo "$$COVERAGE < 80" | bc -l) -eq 1 ]; then \
		echo "FAIL: Coverage $$COVERAGE% is below 80% threshold"; \
		exit 1; \
	fi
	@echo "PASS: Coverage meets 80% threshold"

.PHONY: ci
ci: fmt vet lint test race coverage-check boundary security contracts

.PHONY: ci-extended
ci-extended: ci property golden fuzz-smoke

.PHONY: release-check
release-check: ci integration
	CHECK_STATUS=passed $(MAKE) evidence
	$(MAKE) release-evidence-check

.PHONY: release-check-extended
release-check-extended: ci-extended integration
	CHECK_STATUS=passed $(MAKE) evidence
	$(MAKE) release-evidence-check

.PHONY: release-final-check
release-final-check: release-check
	RELEASE_EVIDENCE_REQUIRE_PASSED=1 RELEASE_EVIDENCE_REQUIRE_CLEAN=1 ./scripts/check_release_evidence.sh

.PHONY: release-preflight
release-preflight:
	./scripts/check_release_preflight.sh "$(VERSION)"
	GOWORK=off VERSION="$(VERSION)" $(MAKE) release-final-check

.PHONY: docker-toolchain-check
docker-toolchain-check:
	./scripts/docker/check_toolchain.sh

.PHONY: docker-build
docker-build:
	DOCKER_BUILDKIT=1 docker buildx build --load --target toolchain \
		--build-arg GO_VERSION="$(GO_VERSION)" \
		--build-arg GOLANGCI_LINT_VERSION="$(GOLANGCI_LINT_VERSION)" \
		--build-arg GOVULNCHECK_VERSION="$(GOVULNCHECK_VERSION)" \
		--tag "$(DOCKER_IMAGE)" .

.PHONY: docker-build-check
docker-build-check:
	$(DOCKER_GATE) build-check

.PHONY: docker-shell
docker-shell: docker-build
	docker run --rm -it --workdir /workspace \
		--volume "$(CURDIR):/workspace" \
		--volume go-build-cache:/root/.cache/go-build \
		--volume go-mod-cache:/go/pkg/mod \
		--env "GOWORK=$${GOWORK:-off}" \
		--env "XLIB_CONTEXT=$${XLIB_CONTEXT:-docker_toolchain}" \
		--env "VERSION=$${VERSION:-}" \
		--env "DOWNSTREAM=$${DOWNSTREAM:-}" \
		--env "XLIB_ENABLE_VULNCHECK=$${XLIB_ENABLE_VULNCHECK:-}" \
		--env "CI=$${CI:-}" \
		--env "GITHUB_ACTIONS=$${GITHUB_ACTIONS:-}" \
		"$(DOCKER_IMAGE)"

.PHONY: docker-ci
docker-ci:
	$(DOCKER_GATE) ci

.PHONY: docker-release-check
docker-release-check:
	$(DOCKER_GATE) release-check

.PHONY: docker-release-final-check
docker-release-final-check:
	$(DOCKER_GATE) release-final-check

.PHONY: docker-goalcli
docker-goalcli:
	$(DOCKER_GATE) goalcli

.PHONY: docker-goalcli-image
docker-goalcli-image: docker-build

.PHONY: docker-goalcli-version
docker-goalcli-version:
	$(DOCKER_GATE) goalcli-version

.PHONY: docker-runtime-check
docker-runtime-check:
	$(DOCKER_GATE) runtime-check

.PHONY: docker-drift-check
docker-drift-check:
	$(DOCKER_GATE) drift-check

.PHONY: docker-contract
docker-contract: docker-toolchain-check docker-build-check docker-runtime-check docker-drift-check
