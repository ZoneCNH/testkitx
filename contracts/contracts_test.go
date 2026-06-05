package contracts

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ZoneCNH/testkitx/pkg/testkitx"
)

type schemaProperty struct {
	Type    string   `json:"type"`
	Enum    []string `json:"enum"`
	Minimum *int     `json:"minimum"`
}

type objectSchema struct {
	Required   []string                  `json:"required"`
	Properties map[string]schemaProperty `json:"properties"`
}

func TestErrorKindContractMatchesPublicConstants(t *testing.T) {
	t.Parallel()

	schema := readSchema(t, "error.schema.json")

	expected := sortedStrings(
		string(testkitx.ErrorKindConfig),
		string(testkitx.ErrorKindValidation),
		string(testkitx.ErrorKindConnection),
		string(testkitx.ErrorKindUnavailable),
		string(testkitx.ErrorKindTimeout),
		string(testkitx.ErrorKindAuth),
		string(testkitx.ErrorKindConflict),
		string(testkitx.ErrorKindRateLimit),
		string(testkitx.ErrorKindInternal),
	)
	actual := sortedStrings(schema.Properties["kind"].Enum...)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("error kind contract drift:\nactual:   %#v\nexpected: %#v", actual, expected)
	}
	requireFields(t, schema.Required, "kind", "op", "message", "retryable")
}

func TestHealthStatusContractMatchesPublicConstants(t *testing.T) {
	t.Parallel()
	schema := readSchema(t, "health.schema.json")

	expected := sortedStrings(
		string(testkitx.HealthHealthy),
		string(testkitx.HealthDegraded),
		string(testkitx.HealthUnhealthy),
	)
	actual := sortedStrings(schema.Properties["status"].Enum...)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("health status contract drift:\nactual:   %#v\nexpected: %#v", actual, expected)
	}
	requireFields(t, schema.Required, "name", "status", "checked_at")
}

func TestConfigContractMatchesPublicConfig(t *testing.T) {
	t.Parallel()
	schema := readSchema(t, "config.schema.json")
	requireFields(t, schema.Required, "name")

	configType := reflect.TypeOf(testkitx.Config{})
	requireSchemaFieldMapsToStructField(t, schema, configType, "name", "Name", "string")
	requireSchemaFieldMapsToStructField(t, schema, configType, "timeout_ms", "Timeout", "integer")
	requireSchemaFieldMapsToStructField(t, schema, configType, "secret", "Secret", "string")

	if timeoutField, ok := configType.FieldByName("Timeout"); !ok || timeoutField.Type != reflect.TypeOf(time.Duration(0)) {
		t.Fatalf("Config.Timeout must remain time.Duration, got %v", timeoutField.Type)
	}
	if minimum := schema.Properties["timeout_ms"].Minimum; minimum == nil || *minimum != 0 {
		t.Fatalf("timeout_ms must define minimum 0, got %#v", minimum)
	}
}

func TestMetricsContractDocumentsPublicConstants(t *testing.T) {
	t.Parallel()
	content, err := os.ReadFile("metrics.md")
	if err != nil {
		t.Fatalf("read metrics contract: %v", err)
	}
	text := string(content)
	for _, metric := range []string{
		testkitx.MetricClientCreatedTotal,
		testkitx.MetricClientClosedTotal,
		testkitx.MetricClientErrorsTotal,
		testkitx.MetricClientHealthStatus,
		testkitx.MetricClientHealthLatencyMS,
		testkitx.MetricClientRequestsTotal,
		testkitx.MetricClientRequestDurationSeconds,
		testkitx.MetricClientRetriesTotal,
		testkitx.MetricClientInflight,
	} {
		if !strings.Contains(text, "`"+metric+"`") {
			t.Fatalf("metrics contract does not document %q", metric)
		}
	}
}

func TestDockerToolchainContractDocumentsEvidenceFields(t *testing.T) {
	t.Parallel()
	schema := readSchema(t, "docker-toolchain.schema.json")

	requireFields(
		t,
		schema.Required,
		"enabled",
		"contract_version",
		"go_version",
		"golangci_lint_version",
		"govulncheck_version",
		"buildkit_required",
		"cache_mounts",
		"validated_by",
		"workflow_run_id",
		"artifact_name",
		"artifact_url",
	)
}

func TestDownstreamAdoptionProofContractIsTestOnly(t *testing.T) {
	t.Parallel()
	schema := readSchema(t, "downstream-adoption-proof.schema.json")

	requireFields(
		t,
		schema.Required,
		"schema_version",
		"source_repo",
		"source_commit",
		"downstream_repo",
		"downstream_commit",
		"mode",
		"adoption_paths",
		"production_import_boundary",
		"gate_outputs",
		"rollback",
	)

	expectedModes := sortedStrings("dry-run", "fixture-only", "test-only")
	actualModes := sortedStrings(schema.Properties["mode"].Enum...)
	if !reflect.DeepEqual(actualModes, expectedModes) {
		t.Fatalf("downstream adoption modes drift:\nactual:   %#v\nexpected: %#v", actualModes, expectedModes)
	}
}

func requireSchemaFieldMapsToStructField(t *testing.T, schema objectSchema, structType reflect.Type, schemaField string, structField string, schemaType string) {
	t.Helper()

	property, ok := schema.Properties[schemaField]
	if !ok {
		t.Fatalf("schema missing property %q", schemaField)
	}
	if property.Type != schemaType {
		t.Fatalf("schema property %q type = %q, want %q", schemaField, property.Type, schemaType)
	}
	if _, ok := structType.FieldByName(structField); !ok {
		t.Fatalf("%s missing field %s required by schema property %q", structType.Name(), structField, schemaField)
	}
}

func readSchema(t *testing.T, path string) objectSchema {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var schema objectSchema
	if err := json.Unmarshal(content, &schema); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return schema
}

func requireFields(t *testing.T, actual []string, expected ...string) {
	t.Helper()
	fields := make(map[string]bool, len(actual))
	for _, field := range actual {
		fields[field] = true
	}
	for _, field := range expected {
		if !fields[field] {
			t.Fatalf("required fields missing %q from %#v", field, actual)
		}
	}
}

func sortedStrings(values ...string) []string {
	copied := append([]string(nil), values...)
	sort.Strings(copied)
	return copied
}
