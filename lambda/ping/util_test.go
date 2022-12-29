package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasTransitionedIntoErrorState(t *testing.T) {
	t.Run("not enough records", func(t *testing.T) {
		records := make([]PingRecord, 0)
		actual := HasTransitionedIntoErrorState(records)
		assert.Equal(t, false, actual)
	})

	t.Run("enough records, no successes", func(t *testing.T) {
		records := []PingRecord{
			{Result: "FAIL"},
			{Result: "FAIL"},
		}
		actual := HasTransitionedIntoErrorState(records)
		assert.Equal(t, true, actual)
	})

	t.Run("not enough failures", func(t *testing.T) {
		records := []PingRecord{
			{Result: "FAIL"},
			{Result: "PASS"},
			{Result: "PASS"},
		}
		actual := HasTransitionedIntoErrorState(records)
		assert.Equal(t, false, actual)
	})

	t.Run("too many failures", func(t *testing.T) {
		records := []PingRecord{
			{Result: "FAIL"},
			{Result: "FAIL"},
			{Result: "FAIL"},
		}
		actual := HasTransitionedIntoErrorState(records)
		assert.Equal(t, false, actual)
	})

	t.Run("has transitioned", func(t *testing.T) {
		records := []PingRecord{
			{Result: "FAIL"},
			{Result: "FAIL"},
			{Result: "PASS"},
		}
		actual := HasTransitionedIntoErrorState(records)
		assert.Equal(t, true, actual)
	})
}
