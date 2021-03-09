package main

import (
	"encoding/json"
	"my-first-telegram-bot/telegram-handler/restclient"
	"my-first-telegram-bot/telegram-handler/utils/mocks"
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

	// t.Run("Successful Request with the mock server", func(t *testing.T) {

	// 	// Arrange
	// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(200)
	// 		if r.Method == http.MethodGet {
	// 			w.Write([]byte("{\"id\": \"96221b11-8a37-4495-baf0-134be4feffc1\", \"text\": \"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.\", \"source\": \"djtech.net\", \"source_url\": \"http://www.djtech.net/humor/useless_facts.htm\", \"language\": \"en\", \"permalink\": \"https://uselessfacts.jsph.pl/96221b11-8a37-4495-baf0-134be4feffc1\"}"))
	// 		}

	// 		if r.Method == http.MethodPost {
	// 			w.Write([]byte("{\"ok\": true,\"result\": {\"message_id\": 26,\"from\": {\"id\": 1025326803,\"is_bot\": true,\"first_name\": \"MyDailyFact\",\"username\": \"majoFFper_bot\"},\"chat\": {\"id\": -255361673,\"title\": \"Pokémons\",\"type\": \"group\",\"all_members_are_administrators\": true},\"date\": 1614894279,\"text\": \"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.\"}}"))
	// 		}

	// 	}))

	// 	restclient.RandomFactsAddress = ts.URL
	// 	restclient.RandomJokesAddress = ts.URL
	// 	restclient.TelegramApi = ts.URL

	// 	defer ts.Close()

	// 	telegramRequest := Update{
	// 		Message: Message{
	// 			Text: "/fact",
	// 			Chat: Chat{
	// 				Id: 1234,
	// 			},
	// 		},
	// 		UpdateId: 1,
	// 	}

	// 	requestBody, err := json.Marshal(telegramRequest)

	// 	tempRequest := events.APIGatewayProxyRequest{
	// 		Body:       string(requestBody),
	// 		Path:       "http://myTelegramWebHookHandler.com/secretToken",
	// 		HTTPMethod: "POST",
	// 	}

	// 	// Act
	// 	response, err := handler(tempRequest)
	// 	if err != nil {
	// 		t.Fatal("Everything should be ok")
	// 	}

	// 	// Assert

	// 	assert.EqualValues(t,
	// 		`{"ok": true,"result": {"message_id": 26,"from": {"id": 1025326803,"is_bot": true,"first_name": "MyDailyFact","username": "majoFFper_bot"},"chat": {"id": -255361673,"title": "Pokémons","type": "group","all_members_are_administrators": true},"date": 1614894279,"text": "To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P."}}`,
	// 		response.Body)
	// })

	t.Run("Issue getting fact from fact api", func(t *testing.T) {

		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		restclient.RandomFactsAddress = ts.URL

		// Act
		_, err := handler(events.APIGatewayProxyRequest{})

		// Assert
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid HTTP response")
		}
	})

	t.Run("Successful Request mocking the rest client", func(t *testing.T) {

		mocks.ReturnGetFact = func() (*restclient.GeneratedFact, error) {

			return &restclient.GeneratedFact{
				ID:        "1",
				Text:      "potato potato",
				Source:    "potato potato",
				SourceURL: "teste123",
				Language:  "en",
				Permalink: "teste124",
			}, nil
		}

		mocks.ReturnPostResponse = func(chatId int, text string) (string, error) {

			escapedJsonContent := "{\"ok\": true,\"result\": {\"message_id\": 26,\"from\": {\"id\": 1025326803,\"is_bot\": true,\"first_name\": \"MyDailyFact\",\"username\": \"majoFFper_bot\"},\"chat\": {\"id\": -255361673,\"title\": \"Pokémons\",\"type\": \"group\",\"all_members_are_administrators\": true},\"date\": 1614894279,\"text\": \"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.\"}}"

			return escapedJsonContent, nil
		}

		telegramRequest := Update{
			Message: Message{
				Text: "/fact",
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

		restclient.MyFactClient = &mocks.MockBaseClient{}

		// Act
		response, err := handler(tempRequest)

		if err != nil {
			t.Fatal("Everything should be ok")
		}

		// Assert

		assert.EqualValues(t,
			`{"ok": true,"result": {"message_id": 26,"from": {"id": 1025326803,"is_bot": true,"first_name": "MyDailyFact","username": "majoFFper_bot"},"chat": {"id": -255361673,"title": "Pokémons","type": "group","all_members_are_administrators": true},"date": 1614894279,"text": "To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P."}}`,
			response.Body)
	})
}
