package compile

import (
	"testing"
)

func TestTopologicalSortNoDeps(t *testing.T) {
	deps := map[string][]string{
		"A": {},
		"B": {},
		"C": {},
	}
	order, err := TopologicalSort(deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 {
		t.Errorf("expected 3 items, got %d: %v", len(order), order)
	}
}

func TestTopologicalSortLinearChain(t *testing.T) {
	// C depends on B, B depends on A. Expected: A, B, C.
	deps := map[string][]string{
		"A": {},
		"B": {"A"},
		"C": {"B"},
	}
	order, err := TopologicalSort(deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := func(name string) int {
		for i, n := range order {
			if n == name {
				return i
			}
		}
		return -1
	}
	if pos("A") > pos("B") || pos("B") > pos("C") {
		t.Errorf("wrong order: %v", order)
	}
}

func TestTopologicalSortExternalRefIgnored(t *testing.T) {
	// B references ENV_VAR which is not in the YAML set — should be ignored.
	deps := map[string][]string{
		"A": {"ENV_VAR"},
		"B": {"A"},
	}
	order, err := TopologicalSort(deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := func(name string) int {
		for i, n := range order {
			if n == name {
				return i
			}
		}
		return -1
	}
	if pos("A") > pos("B") {
		t.Errorf("A must come before B, got: %v", order)
	}
}

func TestTopologicalSortCircular(t *testing.T) {
	deps := map[string][]string{
		"A": {"B"},
		"B": {"A"},
	}
	_, err := TopologicalSort(deps)
	if err == nil {
		t.Error("expected circular reference error")
	}
}

func TestTopologicalSortDeterministic(t *testing.T) {
	deps := map[string][]string{
		"Z": {},
		"A": {},
		"M": {},
	}
	order1, _ := TopologicalSort(deps)
	order2, _ := TopologicalSort(deps)
	for i := range order1 {
		if order1[i] != order2[i] {
			t.Errorf("non-deterministic order: %v vs %v", order1, order2)
		}
	}
}
