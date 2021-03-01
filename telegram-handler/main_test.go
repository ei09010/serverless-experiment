package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// t.Run("Unable to get IP", func(t *testing.T) {
	// 	RandomFactsAddress = "http://127.0.0.1:12345"

	// 	_, err := handler(events.APIGatewayProxyRequest{})
	// 	if err == nil {
	// 		t.Fatal("Error failed to trigger with an invalid request")
	// 	}
	// })

	// t.Run("Non 200 Response", func(t *testing.T) {
	// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(500)
	// 	}))
	// 	defer ts.Close()

	// 	RandomFactsAddress = ts.URL

	// 	_, err := handler(events.APIGatewayProxyRequest{})
	// 	if err != nil && err.Error() != ErrNon200Response.Error() {
	// 		t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
	// 	}
	// })

	// t.Run("Unable decode IP", func(t *testing.T) {
	// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(500)
	// 	}))
	// 	defer ts.Close()

	// 	RandomFactsAddress = ts.URL

	// 	_, err := handler(events.APIGatewayProxyRequest{})
	// 	if err == nil {
	// 		t.Fatal("Error failed to trigger with an invalid HTTP response")
	// 	}
	// })

	t.Run("Successful Request", func(t *testing.T) {

		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			escapedJsonContent := "{\"id\": \"96221b11-8a37-4495-baf0-134be4feffc1\", \"text\": \"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.\", \"source\": \"djtech.net\", \"source_url\": \"http://www.djtech.net/humor/useless_facts.htm\", \"language\": \"en\", \"permalink\": \"https://uselessfacts.jsph.pl/96221b11-8a37-4495-baf0-134be4feffc1\"}"
			w.Write([]byte(escapedJsonContent))
		}))

		defer ts.Close()

		RandomFactsAddress = ts.URL

		telegramRequest := Update{
			Message: Message{
				Text: "hello world",
				Chat: Chat{
					Id: 1234,
				},
			},
			UpdateId: 1,
		}

		requestBody, err := json.Marshal(telegramRequest)

		tempRequest := events.APIGatewayProxyRequest{
			Body:       string(requestBody),
			Path:       "http://myTelegramWebHookHandler.com/secretToken",
			HTTPMethod: "POST",
		}

		// Act
		response, err := handler(tempRequest)
		if err != nil {
			t.Fatal("Everything should be ok")
		}

		// Assert

		assert.EqualValues(t,
			"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.",
			response.Body)
	})
}
