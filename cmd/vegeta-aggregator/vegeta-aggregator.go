package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"vegeta-kubernetes/internal/pkg/aggregator"
	"vegeta-kubernetes/internal/pkg/controllers"
)

var (
	selector = flag.String("selector", "", "The label selector for pods")
	sleep    = flag.Duration("sleep", 5*time.Second, "The sleep period between aggregations")
)

func main() {
	flag.Parse()

	err := aggregator.AggregateData(*selector, *sleep)
	if err != nil {
		log.Fatalln("Failed to aggregate date from pods:", err)
	}

	http.HandleFunc("/metrics", controllers.AggregateMetrics)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
