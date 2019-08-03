package aggregator

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tsenart/vegeta/lib"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	lock             = &sync.Mutex{}
	aggregateMetrics = &vegeta.Metrics{}
)

// GetMetrics return statistics about the load testing
func GetMetrics(w http.ResponseWriter) {
	lock.Lock()
	defer lock.Unlock()
	reporter := vegeta.NewJSONReporter(aggregateMetrics)
	w.WriteHeader(http.StatusOK)
	err := reporter.Report(w)
	if err != nil {
		log.Fatalln("Failed to get metrics:", err)
	}
}

func setData(data []vegeta.Metrics) {
	lock.Lock()
	defer lock.Unlock()

	for _, m := range data {
		*aggregateMetrics = addMetrics(aggregateMetrics, &m)
	}
}

func addMetrics(m1 *vegeta.Metrics, m2 *vegeta.Metrics) vegeta.Metrics {
	duration := m1.Duration + m2.Duration
	requests := m1.Requests + m2.Requests
	totalLatencies := m1.Latencies.Total + m2.Latencies.Total
	meanLatencies := time.Duration(((uint64(m1.Latencies.Mean) * m1.Requests) + (uint64(m2.Latencies.Mean) * m2.Requests)) / (requests))
	maxLatencies := time.Duration(math.Max(float64(m1.Latencies.Max), float64(m2.Latencies.Max)))
	rate := m1.Rate + m2.Rate
	success := m1.Success + m2.Success
	wait := m1.Wait + m2.Wait
	metricsErrors := append(m1.Errors, m2.Errors...)

	return vegeta.Metrics{
		Duration: duration,
		Latencies: vegeta.LatencyMetrics{
			Total: totalLatencies,
			Mean:  meanLatencies,
			Max:   maxLatencies,
		},
		Requests: requests,
		Rate:     rate,
		Success:  success,
		Wait:     wait,
		Errors:   metricsErrors,
	}
}

func AggregateData(selector string, sleep time.Duration) error {
	for {
		start := time.Now()

		if err := loadData(selector); err != nil {
			return errors.Wrap(err, "Error while trying to aggregate metrics")
		}

		latency := time.Now().Sub(start)
		if latency < sleep {
			time.Sleep(sleep - latency)
		}
		fmt.Printf("%v\n", time.Now().Sub(start))
	}

	return nil
}

func loadData(selector string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "Error creating in cluster config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "Error creating clientset config")
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return errors.Wrap(err, "Error getting pods")
	}

	loadbots := []*v1.Pod{}
	for i := range pods.Items {
		pod := &pods.Items[i]
		if pod.Status.PodIP == "" {
			continue
		}
		loadbots = append(loadbots, pod)
	}

	metrics := []vegeta.Metrics{}
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(loadbots))
	for i := range loadbots {
		go func(pod *v1.Pod) error {
			defer wg.Done()

			data, err := clientset.RESTClient().Get().AbsPath("/api/v1/proxy/namespaces/default/pods/" + pod.Name + ":8080/").DoRaw()
			if err != nil {
				return errors.Wrapf(err, "Error proxying to pod: %s\n", pod.Name)
			}

			var metric vegeta.Metrics
			if err := json.Unmarshal(data, &metric); err != nil {
				return errors.Wrap(err, "Failed to unmarshal data to vegeta metrics")
			}

			lock.Lock()
			defer lock.Unlock()

			metrics = append(metrics, metric)
			return nil
		}(loadbots[i])
	}

	wg.Wait()
	setData(metrics)

	reporter := vegeta.NewJSONReporter(aggregateMetrics)
	err = reporter.Report(os.Stdout)
	if err != nil {
		log.Fatalln("Failed to report metrics to stdout:", err)
	}

	return nil
}
