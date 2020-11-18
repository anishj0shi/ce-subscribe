package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anishj0shi/ce-subscriber/pkg/api"
	"net/http"
	"time"
)

type InMemoryDBServiceClient interface {
	SendEventLatency(object api.InMemoryDataObject) error
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

func (i inMemoryDBServiceClient) SendEventLatency(object api.InMemoryDataObject) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, i.url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	res, err := i.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New("Entry not added")
	}
	return nil
}
