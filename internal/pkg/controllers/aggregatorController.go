package controllers

import (
	"net/http"

	"vegeta-kubernetes/internal/pkg/aggregator"
)

func AggregateMetrics(w http.ResponseWriter, r *http.Request) {
	aggregator.GetMetrics(w)
}
