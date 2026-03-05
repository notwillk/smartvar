package validate

import "testing"

func TestValidateEmail(t *testing.T) {
	valid := []string{"user@example.com", "first.last+tag@sub.domain.org"}
	invalid := []string{"", "notanemail", "@example.com", "user@", "user @example.com"}
	runCases(t, "email", valid, invalid)
}

func TestValidateURL(t *testing.T) {
	valid := []string{"http://example.com", "https://api.example.com/path?q=1"}
	invalid := []string{"", "ftp://example.com", "example.com", "not a url"}
	runCases(t, "url", valid, invalid)
}

func TestValidateHostname(t *testing.T) {
	valid := []string{"example.com", "api.example.com", "localhost", "my-host"}
	invalid := []string{"", "192.168.1.1", "label.toolongXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}
	runCases(t, "hostname", valid, invalid)
}

func TestValidateUUID(t *testing.T) {
	valid := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"00000000-0000-0000-0000-000000000000",
	}
	invalid := []string{"", "not-a-uuid", "550e8400e29b41d4a716446655440000"}
	runCases(t, "uuid", valid, invalid)
}

func TestValidateInt(t *testing.T) {
	valid := []string{"0", "42", "-1", "9223372036854775807"}
	invalid := []string{"", "1.5", "abc", "1e5"}
	runCases(t, "int", valid, invalid)
}

func TestValidateBool(t *testing.T) {
	valid := []string{"true", "false", "True", "False", "TRUE", "FALSE", "1", "0"}
	invalid := []string{"", "yes", "no", "2", "truee"}
	runCases(t, "bool", valid, invalid)
}

func TestValidatePort(t *testing.T) {
	valid := []string{"1", "80", "443", "8080", "65535"}
	invalid := []string{"", "0", "65536", "-1", "abc"}
	runCases(t, "port", valid, invalid)
}

func TestValidateUnknownPattern(t *testing.T) {
	if err := Validate("nonexistent", "anything"); err == nil {
		t.Error("expected error for unknown pattern")
	}
}

func runCases(t *testing.T, pattern string, valid, invalid []string) {
	t.Helper()
	for _, v := range valid {
		if err := Validate(pattern, v); err != nil {
			t.Errorf("pattern=%s value=%q: expected valid, got %v", pattern, v, err)
		}
	}
	for _, v := range invalid {
		if err := Validate(pattern, v); err == nil {
			t.Errorf("pattern=%s value=%q: expected invalid, got nil", pattern, v)
		}
	}
}
