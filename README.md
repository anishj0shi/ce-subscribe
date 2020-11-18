# ce-subscriber
Example implementation of a Cloud Event Subscriber. This example listens on
port `8081` and calculates the timestamp difference between the event
triggered time and the event received time. This time difference gives an
overall end to end event latency estimation.

For analysis, we send the latency infomation to a persistence service, where
the results are aggregated into an InMemmory Database. for more information,
also see: https://github.com/anishj0shi/inmemorydb-service
