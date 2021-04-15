package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	host              string
	runForSeconds     int
	requestsPerSecond int
	workers           int
	minLatencyMs      int
	maxLatencyMs      int

	latencyResults = []int64{}
)

type requestParams struct {
	latency    time.Duration
	statusCode int
}

const urlFormat = "%s/?latencyms=%v&statusCode=%v"

func init() {
	flag.StringVar(&host, "host", "", "the host that requests will be made to")
	flag.IntVar(&requestsPerSecond, "rps", 5, "the number of requests to send per second")
	flag.IntVar(&runForSeconds, "for", 10, "the number of seconds to run for")
	flag.IntVar(&workers, "workers", 5, "the number of workers to use")
	flag.IntVar(&minLatencyMs, "min-latency", 0, "the minimum latency in millieconds")
	flag.IntVar(&maxLatencyMs, "max-latency", 0, "the maximum latency in millieconds")

	flag.Parse()

	if host == "" {
		log.Fatal("host arg has not been provided")
	}

	if minLatencyMs > maxLatencyMs {
		log.Fatal("min latency must be less than or equal to max latency")
	}
}

func main() {
	client := &http.Client{}
	ticker := time.NewTicker(1 * time.Second)
	workChan := make(chan *requestParams)
	logChan := make(chan string)

	for i := 0; i < workers; i++ {
		go func(workerId int) {
			for {
				d := <-workChan
				url := fmt.Sprintf(urlFormat, host, d.latency.Milliseconds(), d.statusCode)

				logChan <- fmt.Sprintf("Worker %v -> %v", workerId, url)

				start := time.Now().UnixNano()
				client.Get(url)
				end := time.Now().UnixNano()
				latencyResults = append(latencyResults, end-start)
			}
		}(i)
	}

	go func() {
		for {
			log.Println(<-logChan)
		}
	}()

	start := time.Now()

	go func() {
		for {
			newTime := <-ticker.C

			logChan <- fmt.Sprintf("running for %.0f of %v seconds", newTime.Sub(start).Seconds(), runForSeconds)

			for i := 0; i < requestsPerSecond; i++ {
				workChan <- &requestParams{
					statusCode: 200,
					latency:    time.Duration(randomInRange(minLatencyMs, maxLatencyMs)) * time.Millisecond,
				}
			}
		}
	}()

	<-time.After(time.Duration(runForSeconds+1) * time.Second)

	log.Printf("made %v requests", len(latencyResults))
}

func randomInRange(min, max int) int {
	return rand.Intn(max-min) + min
}
