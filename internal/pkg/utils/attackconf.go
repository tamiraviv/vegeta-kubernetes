package utils

import (
	"net/http"
	"time"
)

type AttackConf struct {
	Url      string
	Method   string
	Headers  http.Header
	Body     []byte
	Rate     int
	Duration time.Duration
	Workers  uint64
}
