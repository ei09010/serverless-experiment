package main

import (
	"encoding/json"
	"my-first-telegram-bot/telegram-handler/restclient"
	"my-first-telegram-bot/telegram-handler/utils/mocks"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	t.Run("Failed Post Response", func(t *testing.T) {

		expectedTelegramResponse := "issue in telegram"

		mocks.ReturnGetJoke = func() (*restclient.GeneratedJoke, error) {

			return &restclient.GeneratedJoke{
				Type: "1",
				Value: restclient.JokeValue{
					ID:         1,
					Joke:       "",
					Categories: []string{"1", "2"},
				},
			}, nil
		}

		mocks.ReturnPostResponse = func(chatId int, text string) (string, error) {
			return expectedTelegramResponse, ErrNon200Response
		}

		telegramRequest := Update{
			Message: Message{
				Text: "/joke",
				Chat: Chat{
					Id: 1234,
				},
			},
			UpdateId: 1,
		}

		requestBody, err := json.Marshal(telegramRequest)

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		tempRequest := events.APIGatewayProxyRequest{
			Body:       string(requestBody),
			Path:       "http://myTelegramWebHookHandler.com/secretToken",
			HTTPMethod: "POST",
		}

		myMockClient := &mocks.MockBaseClient{}

		restclient.MyJokeClient = myMockClient

		restclient.MyTelegramClient = myMockClient

		// Act
		response, err := handler(tempRequest)

		// Assert

		assert.Equal(t, 1, myMockClient.ReturnGetJokeCallCount)

		assert.Equal(t, 0, myMockClient.ReturnGetFactCallCount)

		assert.Equal(t, 1, myMockClient.ReturnPostResponseCallCount)

		assert.EqualValues(t,
			expectedTelegramResponse,
			response.Body)
	})

	t.Run("Failed Fact Request", func(t *testing.T) {

		mocks.ReturnGetFact = func() (*restclient.GeneratedFact, error) {

			return nil, ErrNon200Response
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

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		tempRequest := events.APIGatewayProxyRequest{
			Body:       string(requestBody),
			Path:       "http://myTelegramWebHookHandler.com/secretToken",
			HTTPMethod: "POST",
		}

		myMockClient := &mocks.MockBaseClient{}

		restclient.MyFactClient = myMockClient

		restclient.MyTelegramClient = myMockClient

		// Act
		response, err := handler(tempRequest)

		// Assert

		assert.Equal(t, 0, myMockClient.ReturnGetJokeCallCount)

		assert.Equal(t, 1, myMockClient.ReturnGetFactCallCount)

		assert.Equal(t, 0, myMockClient.ReturnPostResponseCallCount)

		assert.EqualValues(t,
			"Error calling fact generation api",
			response.Body)
	})

	t.Run("Failed Joke Request", func(t *testing.T) {

		mocks.ReturnGetJoke = func() (*restclient.GeneratedJoke, error) {

			return nil, ErrNon200Response
		}

		telegramRequest := Update{
			Message: Message{
				Text: "/joke",
				Chat: Chat{
					Id: 1234,
				},
			},
			UpdateId: 1,
		}

		requestBody, err := json.Marshal(telegramRequest)

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		tempRequest := events.APIGatewayProxyRequest{
			Body:       string(requestBody),
			Path:       "http://myTelegramWebHookHandler.com/secretToken",
			HTTPMethod: "POST",
		}

		myMockClient := &mocks.MockBaseClient{}

		restclient.MyJokeClient = myMockClient

		restclient.MyTelegramClient = myMockClient

		// Act
		response, err := handler(tempRequest)

		// Assert

		assert.Equal(t, 1, myMockClient.ReturnGetJokeCallCount)

		assert.Equal(t, 0, myMockClient.ReturnGetFactCallCount)

		assert.Equal(t, 0, myMockClient.ReturnPostResponseCallCount)

		assert.EqualValues(t,
			"Error calling joke generation api",
			response.Body)
	})

	t.Run("Successful Joke Request", func(t *testing.T) {

		mocks.ReturnGetJoke = func() (*restclient.GeneratedJoke, error) {

			return &restclient.GeneratedJoke{
				Type: "1",
				Value: restclient.JokeValue{
					ID:         1,
					Joke:       "",
					Categories: []string{"1", "2"},
				},
			}, nil
		}

		mocks.ReturnPostResponse = func(chatId int, text string) (string, error) {

			escapedJsonContent := "{\"ok\": true,\"result\": {\"message_id\": 26,\"from\": {\"id\": 1025326803,\"is_bot\": true,\"first_name\": \"MyDailyFact\",\"username\": \"majoFFper_bot\"},\"chat\": {\"id\": -255361673,\"title\": \"Pokémons\",\"type\": \"group\",\"all_members_are_administrators\": true},\"date\": 1614894279,\"text\": \"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.\"}}"

			return escapedJsonContent, nil
		}

		telegramRequest := Update{
			Message: Message{
				Text: "/joke",
				Chat: Chat{
					Id: 1234,
				},
			},
			UpdateId: 1,
		}

		requestBody, err := json.Marshal(telegramRequest)

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		tempRequest := events.APIGatewayProxyRequest{
			Body:       string(requestBody),
			Path:       "http://myTelegramWebHookHandler.com/secretToken",
			HTTPMethod: "POST",
		}

		myMockClient := &mocks.MockBaseClient{}

		restclient.MyJokeClient = myMockClient

		restclient.MyTelegramClient = myMockClient

		// Act
		response, err := handler(tempRequest)

		// Assert

		assert.Equal(t, 1, myMockClient.ReturnGetJokeCallCount)

		assert.Equal(t, 0, myMockClient.ReturnGetFactCallCount)

		assert.Equal(t, 1, myMockClient.ReturnPostResponseCallCount)

		assert.EqualValues(t,
			`{"ok": true,"result": {"message_id": 26,"from": {"id": 1025326803,"is_bot": true,"first_name": "MyDailyFact","username": "majoFFper_bot"},"chat": {"id": -255361673,"title": "Pokémons","type": "group","all_members_are_administrators": true},"date": 1614894279,"text": "To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P."}}`,
			response.Body)
	})

	t.Run("Successful Fact Request", func(t *testing.T) {

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

		myMockClient := &mocks.MockBaseClient{}

		restclient.MyFactClient = myMockClient

		restclient.MyTelegramClient = myMockClient

		// Act
		response, err := handler(tempRequest)

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		// Assert

		assert.Equal(t, 1, myMockClient.ReturnPostResponseCallCount)

		assert.Equal(t, 1, myMockClient.ReturnGetFactCallCount)

		assert.Equal(t, 0, myMockClient.ReturnGetJokeCallCount)

		assert.EqualValues(t,
			`{"ok": true,"result": {"message_id": 26,"from": {"id": 1025326803,"is_bot": true,"first_name": "MyDailyFact","username": "majoFFper_bot"},"chat": {"id": -255361673,"title": "Pokémons","type": "group","all_members_are_administrators": true},"date": 1614894279,"text": "To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P."}}`,
			response.Body)
	})
}
