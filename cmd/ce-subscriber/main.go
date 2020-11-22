package main

import (
	"context"
	"github.com/anishj0shi/ce-subscriber/pkg/client"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"log"
	"net/http"
	"os"
)

const (
	INMEMORY_DB_SERVICE_URL_KEY = "INMEMORY_DB_SERVICE_URL"
)

var inMemoryDBServiceURL string
var dbClient client.InMemoryDBServiceClient

func main() {
	log.Println("Starting event Receiver...")
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
	dbClient.SendEventLatency(event)
}
