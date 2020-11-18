package main

import (
	"context"
	"github.com/anishj0shi/ce-subscriber/pkg/api"
	"github.com/anishj0shi/ce-subscriber/pkg/client"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/exp/rand"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	INMEMORY_DB_SERVICE_URL_KEY = "INMEMORY_DB_SERVICE_URL"
)

var inMemoryDBServiceURL string
var dbClient client.InMemoryDBServiceClient

func main() {
	inMemoryDBServiceURL = os.Getenv(INMEMORY_DB_SERVICE_URL_KEY)
	if inMemoryDBServiceURL == "" {
		panic("No Reference to InMemoryDBService. " +
			"Ensure InMemoryDBService service is deployed. " +
			"ref: https://github.com/anishj0shi/inmemorydb-service")
	}
	dbClient = client.NewInMemoryDBServiceClient(inMemoryDBServiceURL)

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatal(err.Error())
	}
	handler, err := cloudevents.NewHTTPReceiveHandler(context.Background(), p, receiveFn)
	if err != nil {
		panic(err.Error())
	}
	if err := http.ListenAndServe(":8081", handler); err != nil {
		panic(err.Error())
	}
}

func receiveFn(ctx context.Context, event event.Event) {
	log.Printf("Received event...")
	eventData := &api.EventData{}
	err := event.DataAs(eventData)
	if err != nil {
		log.Print("unable to deserialise data object")
	}
	e2eLatency := getLatencyDifference(time.Now().Unix(), eventData.Timestamp)
	obj := api.InMemoryDataObject{
		ID:         rand.Int(),
		EventId:    event.ID(),
		E2ELatency: e2eLatency,
		EventType:  event.Type(),
	}

	err = dbClient.SendEventLatency(obj)
	if err != nil {
		log.Print(err)
	}
}

func getLatencyDifference(now, timestamp int64) int64 {
	return now - timestamp
}
