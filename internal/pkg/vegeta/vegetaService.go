package vegeta

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"vegeta-kubernetes/internal/pkg/utils"

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
	err := reporter.Report(w)
	if err != nil {
		log.Fatalln("Failed to get metrics:", err)
	}
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
func Attack(ac utils.AttackConf) {
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: ac.Method,
		URL:    ac.Url,
		Header: ac.Headers,
		Body:   ac.Body,
	})

	attacker := vegeta.NewAttacker(vegeta.Workers(ac.Workers))

	fmt.Printf("Starting to attack %s!\n", ac.Url)
	for {
		//metrics := &vegeta.Metrics{}
		results := vegeta.Results{}
		vegetaRate := vegeta.Rate{ac.Rate, 1 * time.Second}
		for res := range attacker.Attack(targeter, vegetaRate, ac.Duration, "attack") {
			fmt.Printf("Got the following result: %+v\n", res)
			results.Add(res)
		}

		results.Close()
		setMetrics(results)

		reporter := vegeta.NewJSONReporter(globalMetrics)
		err := reporter.Report(os.Stdout)
		if err != nil {
			log.Fatalln("Failed to report metrics to stdout:", err)
		}
	}
}
