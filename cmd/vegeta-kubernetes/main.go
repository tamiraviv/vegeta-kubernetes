package main

import (
	"flag"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	route "vegeta-kubernetes/internal/pkg/routers"
	"vegeta-kubernetes/internal/pkg/vegeta"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	host     = flag.String("host", "localhost", "The host to load test")
	port     = flag.Int("port", 80, "The port to load test")
	rate     = flag.Int("rate", 1000, "The QPS to send")
	duration = flag.Duration("duration", 1*time.Second, "The duration of the load test")
	workers  = flag.Uint64("workers", 100, "The number of workers to use")
	endPoint = flag.String("endPoint", "/", "The end point to load test")
)

func main() {
	flag.Parse()

	/*serviceIP, err := convertHostToIP(*host)
	if err != nil {
		log.Fatalln("Error while trying to convert host to ip: %s", err)
	}*/

	url, err := buildURL(*host+":"+strconv.Itoa(*port), *endPoint)
	if err != nil {
		log.Fatalln(err)
	}

	r := route.NewRouter()

	go vegeta.Attack(url, *rate, *duration, *workers)

	http.ListenAndServe(":8080", r)
}

func buildURL(host string, endPoint string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", errors.Wrapf(err, "Error while trying to parse host %s to url", host)
	}
	u.Path = path.Join(u.Path, "bar/")
	return u.String(), nil
}

/*func convertHostToIP(host string) (string,error) {
	var serviceIP string
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		ipv4 := ip.To4()
		if ipv4 != nil {
			serviceIP = ipv4.String()
			break
		}
	}

	if len(serviceIP) == 0 {
		return "", errors.New("Could not")
	}

	return serviceIP, nil
}*/
