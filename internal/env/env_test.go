package env

import (
	"strings"
	"testing"
)

func TestParseEnvLines(t *testing.T) {
	input := `
# comment
KEY1=value1
KEY2=value with spaces
KEY3=has=equals=in=value
  KEY4=trimmed line

`
	m, err := ParseEnvLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := map[string]string{
		"KEY1": "value1",
		"KEY2": "value with spaces",
		"KEY3": "has=equals=in=value",
		"KEY4": "trimmed line",
	}
	for k, want := range cases {
		if got := m[k]; got != want {
			t.Errorf("%s: got %q, want %q", k, got, want)
		}
	}
	if len(m) != 4 {
		t.Errorf("expected 4 keys, got %d: %v", len(m), m)
	}
}

func TestParseEnvLinesEmpty(t *testing.T) {
	m, err := ParseEnvLines(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestMerge(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "overridden", "C": "3"}
	got := Merge(a, b)

	if got["A"] != "1" {
		t.Errorf("A: got %q, want 1", got["A"])
	}
	if got["B"] != "overridden" {
		t.Errorf("B: got %q, want overridden", got["B"])
	}
	if got["C"] != "3" {
		t.Errorf("C: got %q, want 3", got["C"])
	}
}

func TestMergeDoesNotMutateInputs(t *testing.T) {
	a := map[string]string{"X": "original"}
	b := map[string]string{"X": "new"}
	Merge(a, b)
	if a["X"] != "original" {
		t.Error("Merge mutated input map a")
	}
}
