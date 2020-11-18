package main

import (
	"fmt"
	"github.com/anishj0shi/ce-subscriber/pkg/api"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/motemen/go-loghttp"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	transport := &loghttp.Transport{
		LogRequest: func(req *http.Request) {
			log.Printf("[%p] %s %s", req, req.Method, req.URL)
		},
		LogResponse: func(resp *http.Response) {
			log.Printf("[%p] %d %s %v", resp.Request, resp.StatusCode, resp.Request.URL, resp.Body)
		},
	}

	client := &http.Client{
		Transport: transport,
	}
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 4 * time.Second
	targeter := getVegetaTarget()
	attacker := vegeta.NewAttacker(vegeta.Client(client))

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	vegeta.NewTextReporter(&metrics).Report(os.Stdout)

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func getVegetaTarget() vegeta.Targeter {
	return func(t *vegeta.Target) error {
		t.Method = http.MethodPost
		t.URL = "http://localhost:8081"
		t.Header = map[string][]string{
			"Content-Type": {"application/cloudevents+json"},
		}

		event := cloudevents.NewEvent()
		event.SetID(string(rand.Int()))
		event.SetSource("pingsource")
		event.SetType("ping.received")
		event.SetData("application/json", &api.EventData{
			Timestamp: time.Now().Unix(),
			UUID:      uuid.New().String(),
		})
		jsonStr, err := event.MarshalJSON()
		if err != nil {
			return err
		}
		t.Body = jsonStr
		return nil
	}
}
