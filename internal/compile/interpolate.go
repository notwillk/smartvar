package compile

import "regexp"

var varRefRe = regexp.MustCompile(`\$\{([^}]+)\}`)

// Interpolate replaces ${VAR} placeholders in s using the provided env map.
// References with no matching key are left as-is.
func Interpolate(s string, env map[string]string) string {
	return varRefRe.ReplaceAllStringFunc(s, func(match string) string {
		m := varRefRe.FindStringSubmatch(match)
		if val, ok := env[m[1]]; ok {
			return val
		}
		return match
	})
}

// ExtractRefs returns all variable names referenced via ${...} in s.
func ExtractRefs(s string) []string {
	matches := varRefRe.FindAllStringSubmatch(s, -1)
	refs := make([]string, 0, len(matches))
	for _, m := range matches {
		refs = append(refs, m[1])
	}
	return refs
}
