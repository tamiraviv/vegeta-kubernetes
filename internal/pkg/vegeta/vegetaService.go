package vegeta

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tsenart/vegeta/lib"
)

var (
	mutex         = &sync.Mutex{}
	globalMetrics = &vegeta.Metrics{}
)

// GetMetrics return statistics about the load testing
func GetMetrics(w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()
	reporter := vegeta.NewJSONReporter(globalMetrics)
	w.WriteHeader(http.StatusOK)
	reporter.Report(w)
}

func setMetrics(results vegeta.Results) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, res := range results {
		globalMetrics.Add(&res)
	}
	globalMetrics.Close()
}

// Attack start the load testing process
func Attack(url string, rate int, duration time.Duration, workers uint64) {
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://" + url,
	})

	attacker := vegeta.NewAttacker(vegeta.Workers(workers))

	fmt.Printf("Starting to attack %s!\n", url)
	for {
		//metrics := &vegeta.Metrics{}
		results := vegeta.Results{}
		vegetaRate := vegeta.Rate{rate, 1 * time.Second}
		for res := range attacker.Attack(targeter, vegetaRate, duration, "attack") {
			results.Add(res)
		}

		results.Close()
		setMetrics(results)

		reporter := vegeta.NewJSONReporter(globalMetrics)
		reporter.Report(os.Stdout)
	}
}
