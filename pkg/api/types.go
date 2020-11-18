package api

type EventData struct {
	Timestamp int64  `json:"timestamp"`
	UUID      string `json:"uuid"`
}

type InMemoryDataObject struct {
	ID         int
	EventId    string
	E2ELatency int64
	EventType  string
}
