package fixture_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestLoadReadsFixtureBytes(t *testing.T) {
	t.Parallel()
	data := fixture.Load(t, "sample.json")
	if len(data) == 0 {
		t.Fatal("expected non-empty fixture data")
	}
}

func TestLoadJSONDeserializesFixture(t *testing.T) {
	t.Parallel()
	var v struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	fixture.LoadJSON(t, "sample.json", &v)
	if v.Name != "test" || v.Value != 42 {
		t.Fatalf("unexpected fixture value: %+v", v)
	}
}
