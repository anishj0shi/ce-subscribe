package client

import (
	"bytes"
	"encoding/json"
	"github.com/anishj0shi/ce-subscriber/pkg/api"
	v2 "github.com/cloudevents/sdk-go/v2"
	"golang.org/x/exp/rand"
	"log"
	"net/http"
	"time"
)

type InMemoryDBServiceClient interface {
	SendEventLatency(event v2.Event)
}

func NewInMemoryDBServiceClient(inmemorydbServiceUrl string) InMemoryDBServiceClient {
	client := &http.Client{
		Timeout: 40 * time.Second,
	}
	return &inMemoryDBServiceClient{
		client: client,
		url:    inmemorydbServiceUrl,
	}
}

type inMemoryDBServiceClient struct {
	client *http.Client
	url    string
}

func (i inMemoryDBServiceClient) SendEventLatency(event v2.Event) {
	currentTime := time.Now().UTC()
	log.Printf("Timestamp: %d, Received event...%v", currentTime, event)

	eventData := &api.EventData{}
	err := event.DataAs(eventData)
	if err != nil {
		log.Print("unable to deserialise data object")
	}
	eventTimeStamp := time.Unix(0, eventData.Timestamp)
	e2eLatency := getLatencyDifference(currentTime, eventTimeStamp)
	obj := api.InMemoryDataObject{
		ID:         rand.Int(),
		EventId:    event.ID(),
		E2ELatency: e2eLatency,
		EventType:  event.Type(),
	}
	b, err := json.Marshal(obj)
	if err != nil {
		log.Print(err)
	}
	req, err := http.NewRequest(http.MethodPost, i.url, bytes.NewBuffer(b))
	if err != nil {
		log.Print(err)
	}
	res, err := i.client.Do(req)
	if err != nil {
		log.Print(err)
	}
	if res.StatusCode != http.StatusCreated {
		log.Print("Entry not added")
	}

}

func getLatencyDifference(now, timestamp time.Time) int64 {
	return now.Sub(timestamp).Milliseconds()
}
