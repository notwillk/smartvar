package validate

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	uuidRe  = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// Validate checks value against the named pattern. Returns nil if valid.
func Validate(pattern, value string) error {
	switch pattern {
	case "email":
		if !emailRe.MatchString(value) {
			return fmt.Errorf("not a valid email address")
		}
	case "url":
		u, err := url.ParseRequestURI(value)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return fmt.Errorf("not a valid HTTP/HTTPS URL")
		}
	case "hostname":
		if net.ParseIP(value) != nil {
			return fmt.Errorf("not a valid hostname (got IP address)")
		}
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("not a valid DNS hostname")
		}
		for _, label := range strings.Split(value, ".") {
			if len(label) == 0 || len(label) > 63 {
				return fmt.Errorf("not a valid DNS hostname")
			}
		}
	case "uuid":
		if !uuidRe.MatchString(value) {
			return fmt.Errorf("not a valid UUID")
		}
	case "int":
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("not a valid integer")
		}
	case "bool":
		lower := strings.ToLower(value)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Errorf("not a valid boolean (expected true/false/1/0)")
		}
	case "port":
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 || n > 65535 {
			return fmt.Errorf("not a valid TCP port (1-65535)")
		}
	default:
		return fmt.Errorf("unknown pattern: %s", pattern)
	}
	return nil
}

// Patterns returns the list of built-in pattern names.
func Patterns() []string {
	return []string{"email", "url", "hostname", "uuid", "int", "bool", "port"}
}
