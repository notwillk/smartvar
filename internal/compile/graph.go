package compile

import (
	"fmt"
	"sort"
)

// TopologicalSort returns names in an order where every dependency comes before
// the variable that depends on it. deps maps each name to the list of names it
// depends on (only YAML-defined names appear as dependency targets; env-only
// refs are resolved at interpolation time and need no ordering).
//
// Returns an error if a circular reference is detected.
func TopologicalSort(deps map[string][]string) ([]string, error) {
	// inDegree[A] = number of YAML deps A still needs resolved before it can run.
	inDegree := make(map[string]int, len(deps))
	// revEdges[B] = vars that depend on B (i.e. must be processed after B).
	revEdges := make(map[string][]string, len(deps))

	for name := range deps {
		inDegree[name] = 0
	}
	for name, depList := range deps {
		for _, dep := range depList {
			if _, inSet := deps[dep]; !inSet {
				continue // dep not in YAML; resolved from env at interpolation time
			}
			inDegree[name]++
			revEdges[dep] = append(revEdges[dep], name)
		}
	}

	// Seed queue with zero-in-degree nodes (deterministic order).
	queue := make([]string, 0, len(deps))
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue)

	result := make([]string, 0, len(deps))
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		next := append([]string(nil), revEdges[node]...)
		sort.Strings(next)
		for _, dep := range next {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = insertSorted(queue, dep)
			}
		}
	}

	if len(result) != len(deps) {
		return nil, fmt.Errorf("circular reference detected")
	}
	return result, nil
}

func insertSorted(s []string, v string) []string {
	i := sort.SearchStrings(s, v)
	s = append(s, "")
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}
