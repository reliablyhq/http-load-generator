package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	host              string
	runForSeconds     int
	requestsPerSecond int
	workers           int
	minLatencyMs      int
	maxLatencyMs      int
	statusCodeFetcher func() int

	latencyResults = []int64{}
)

type requestParams struct {
	latency    time.Duration
	statusCode int
}

const urlFormat = "%s/?latency=%v&statuscode=%v"

func init() {
	var statusCode int
	var statusCodeOptionsString string

	flag.StringVar(&host, "host", "", "the host that requests will be made to")
	flag.IntVar(&requestsPerSecond, "rps", 5, "the number of requests to send per second")
	flag.IntVar(&runForSeconds, "for", 10, "the number of seconds to run for")
	flag.IntVar(&workers, "workers", 5, "the number of workers to use")
	flag.IntVar(&minLatencyMs, "min-latency", 0, "the minimum latency in millieconds")
	flag.IntVar(&maxLatencyMs, "max-latency", 0, "the maximum latency in millieconds")
	flag.IntVar(&statusCode, "status-code", -1, "the status code you want to return")
	flag.StringVar(&statusCodeOptionsString, "status-codes", "", "a collection of status codes to return randomly")

	flag.Parse()

	if host == "" {
		log.Fatal("host arg has not been provided")
	}

	if minLatencyMs > maxLatencyMs {
		log.Fatal("min latency must be less than or equal to max latency")
	}

	if statusCode > -1 {
		statusCodeFetcher = func() int {
			return statusCode
		}
	} else if statusCodeOptionsString != "" {
		parts := strings.Split(statusCodeOptionsString, ",")
		ints := []int{}
		for index, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				ints = append(ints, i)
			} else {
				log.Fatalf("status code option %v is not a valid int", index)
			}
		}

		statusCodeFetcher = func() int {
			return ints[rand.Intn(len(ints))]
		}
	} else {
		statusCodeFetcher = func() int {
			return 200
		}
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
				url := fmt.Sprintf(urlFormat, host, d.latency.String(), d.statusCode)

				logChan <- fmt.Sprintf("Worker %v -> %v", workerId, url)

				start := time.Now().UnixNano()
				if _, err := client.Get(url); err != nil {
					log.Print(err)
					continue
				}
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
					statusCode: statusCodeFetcher(),
					latency:    time.Duration(randomInRange(minLatencyMs, maxLatencyMs)) * time.Millisecond,
				}
			}
		}
	}()

	<-time.After(time.Duration(runForSeconds+1) * time.Second)

	log.Printf("made %v requests", len(latencyResults))
}

func randomInRange(min, max int) int {
	if min == max {
		return min
	}

	return rand.Intn(max-min) + min
}
