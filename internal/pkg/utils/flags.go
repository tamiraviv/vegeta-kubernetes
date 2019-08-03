package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Headers struct{ http.Header }

func (h Headers) String() string {
	buf := &bytes.Buffer{}
	if err := h.Write(buf); err != nil {
		return ""
	}
	return buf.String()
}

// Set implements the flag.Value interface for a map of HTTP Headers.
func (h Headers) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("header '%s' has a wrong format", value)
	}
	key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if key == "" || val == "" {
		return fmt.Errorf("header '%s' has a wrong format", value)
	}
	// Add key/value directly to the http.Header (map[string][]string).
	// http.Header.Add() canonicalizes keys but vegeta is used
	// to test systems that require case-sensitive headers.
	h.Header[key] = append(h.Header[key], val)
	return nil
}
