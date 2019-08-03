package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"
	"vegeta-kubernetes/internal/pkg/controllers"
	"vegeta-kubernetes/internal/pkg/utils"

	"vegeta-kubernetes/internal/pkg/customtests"
	"vegeta-kubernetes/internal/pkg/vegeta"
)

var (
	target   = flag.String("url", "http://localhost:80/", "The host to load test")
	headers  = utils.Headers{http.Header{}}
	body     = flag.String("body", "", "The body for the request")
	rate     = flag.Int("rate", 1000, "The QPS to send")
	duration = flag.Duration("duration", 1*time.Second, "The duration of the load test")
	workers  = flag.Uint64("workers", 100, "The number of workers to use")
	method   = flag.String("method", "GET", "The method of the request")

	// optional run custom test. When running custom test, method flag ignored.
	testName = flag.String("testName", "", "The custom test name to run")
)

func main() {
	flag.Var(&headers, "headers", "The headers for the request")
	flag.Parse()

	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		log.Fatalln("Url must contain prefix 'http' or 'https'")
	}

	var requestBody []byte
	if *body == "" {
		requestBody = nil
	} else {
		var err error
		requestBody, err = json.Marshal(*body)
		if err != nil {
			log.Fatalf("Failed to marshal body: %s, Error: %s", *body, err)
		}
	}

	ac := utils.AttackConf{
		Url:      *target,
		Method:   *method,
		Headers:  headers.Header,
		Body:     requestBody,
		Rate:     *rate,
		Duration: *duration,
		Workers:  *workers,
	}

	if *testName == "" {
		go vegeta.Attack(ac)
	} else {
		go customtests.Run(*testName, ac)
	}

	http.HandleFunc("/metrics", controllers.GetMetrics)
	http.ListenAndServe(":8080", nil)
}
