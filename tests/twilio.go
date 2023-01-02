package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func getTwilioMessage(t *testing.T, msgText string) *twilioApi.ApiV2010Message {
	var msg *twilioApi.ApiV2010Message
	client := twilio.NewRestClient()
	params := twilioApi.ListMessageParams{}
	params.SetDateSentAfter(time.Now().Add(-1 * time.Minute))
	params.SetLimit(1)

	for attempt := 0; attempt < 3; attempt++ {
		messages, err := client.Api.ListMessage(&params)
		require.NoError(t, err)
		if len(messages) > 0 && messages[0].Body != nil && strings.Contains(*messages[0].Body, msgText) {
			msg = &messages[0]
			break
		}
		time.Sleep(10 * time.Second)
	}

	return msg
}
