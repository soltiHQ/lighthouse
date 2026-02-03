package domain

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func validateStringNotEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w (%s is empty)", ErrFieldEmpty, field)
	}
	return nil
}

func validateURL(field, value string) error {
	if err := validateStringNotEmpty(field, value); err != nil {
		return err
	}

	if !strings.Contains(value, "://") {
		return fmt.Errorf("%w (%s must include scheme, e.g., 'http://example.com')", ErrInvalidURL, field)
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("%w (%s invalid format: %q)", ErrInvalidURL, field, value)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("%w (%s unsupported scheme %q, only http/https allowed)", ErrInvalidURL, field, parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("%w (%s missing host)", ErrInvalidURL, field)
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return fmt.Errorf("%w (%s missing hostname)", ErrInvalidURL, field)
	}
	if portStr := parsed.Port(); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("%w (%s invalid port %q)", ErrInvalidURL, field, portStr)
		}
		if err = validatePortForScheme(field, parsed.Scheme, port); err != nil {
			return err
		}
	}

	if parsed.User != nil {
		return fmt.Errorf("%w (%s contains credentials in URL, use headers instead)", ErrInvalidURL, field)
	}
	return nil
}

func validatePortForScheme(field, scheme string, port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("%w (%s port %d out of valid range 1-65535)", ErrInvalidURL, field, port)
	}
	return nil
}
