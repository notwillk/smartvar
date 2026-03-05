package compile

import (
	"reflect"
	"sort"
	"testing"
)

func TestInterpolate(t *testing.T) {
	env := map[string]string{
		"HOST": "example.com",
		"PORT": "8080",
	}
	cases := []struct {
		in, want string
	}{
		{"no refs here", "no refs here"},
		{"${HOST}", "example.com"},
		{"https://${HOST}:${PORT}", "https://example.com:8080"},
		{"${UNDEFINED}", "${UNDEFINED}"},
		{"${HOST}/${UNDEFINED}", "example.com/${UNDEFINED}"},
		{"", ""},
	}
	for _, c := range cases {
		if got := Interpolate(c.in, env); got != c.want {
			t.Errorf("Interpolate(%q): got %q, want %q", c.in, got, c.want)
		}
	}
}

func TestExtractRefs(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"no refs", nil},
		{"${A}", []string{"A"}},
		{"${A} and ${B}", []string{"A", "B"}},
		{"https://${HOST}:${PORT}/path", []string{"HOST", "PORT"}},
		{"", nil},
	}
	for _, c := range cases {
		got := ExtractRefs(c.in)
		if len(got) == 0 && len(c.want) == 0 {
			continue
		}
		sort.Strings(got)
		sort.Strings(c.want)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ExtractRefs(%q): got %v, want %v", c.in, got, c.want)
		}
	}
}
