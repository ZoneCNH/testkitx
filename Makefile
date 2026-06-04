.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: race
race:
	go test -race ./...

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
	@go test ./... -coverprofile=coverage.out -covermode=atomic -timeout 180s
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
