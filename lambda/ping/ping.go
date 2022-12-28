package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type PingRecord struct {
	Endpoint       string `dynamodbav:"Endpoint"`
	ExpirationTime int64  `dynamodbav:"ExpirationTime"`
	Result         string `dynamodbav:"Result"`
	Timestamp      int64  `dynamodbav:"Timestamp"`
}

var httpClient = http.Client{Timeout: 10 * time.Second}

func PingSingleEndpoint(endpoint string) PingRecord {
	log.Println("sending request to", endpoint)
	resp, err := httpClient.Get(endpoint)

	now := time.Now()
	expiration := now.Add(24 * time.Hour)
	record := PingRecord{Endpoint: endpoint, ExpirationTime: expiration.Unix(), Timestamp: now.Unix()}

	if err != nil {
		log.Printf("received error while sending request to %s: %s\n", endpoint, err)
		record.Result = "FAIL"
	} else if resp.StatusCode >= 400 {
		log.Printf("received error status code from request to %s: %d\n", endpoint, resp.StatusCode)
		record.Result = "FAIL"
	} else {
		log.Println("received successful response from", endpoint)
		record.Result = "PASS"
	}

	return record
}

func PingAllEndpoints(endpoints []string) []PingRecord {
	var wg sync.WaitGroup
	wg.Add(len(endpoints))
	results := make([]PingRecord, len(endpoints))

	for i, endpoint := range endpoints {
		go func(i int, endpoint string) {
			defer wg.Done()
			results[i] = PingSingleEndpoint(endpoint)
		}(i, endpoint)
	}
	wg.Wait()

	return results
}
