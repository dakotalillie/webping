package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingSingleEndpoint(t *testing.T) {
	t.Run("no protocol", func(t *testing.T) {
		record := PingSingleEndpoint("google.com")
		assert.Equal(t, "FAIL", record.Result)
	})

	t.Run("404 endpoint", func(t *testing.T) {
		record := PingSingleEndpoint("https://google.com/404")
		assert.Equal(t, "FAIL", record.Result)
	})

	t.Run("success", func(t *testing.T) {
		record := PingSingleEndpoint("https://google.com")
		assert.Equal(t, "PASS", record.Result)
	})
}

func TestPingAllEndpoints(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		records := PingAllEndpoints([]string{"https://google.com", "https://apple.com"})
		for _, record := range records {
			assert.Equal(t, "PASS", record.Result, fmt.Sprintf("received error response for %s", record.Endpoint))
		}
	})
}
