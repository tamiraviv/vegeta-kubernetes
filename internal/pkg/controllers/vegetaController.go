package controllers

import (
	"net/http"
	"vegeta-kubernetes/internal/pkg/vegeta"
)

// GetMetrics return statistics about the load testing
func GetMetrics(w http.ResponseWriter, r *http.Request) {
	vegeta.GetMetrics(w)
}
