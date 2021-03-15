package restclient

import (
	"bytes"
	"errors"
	"io/ioutil"
	"my-first-telegram-bot/telegram-handler/utils/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailedFactRequest(t *testing.T) {

	t.Run("Error fact request", func(t *testing.T) {

		expectedId := ""
		expectedText := ""
		expectedSourceUrl := ""
		expectedLanguage := ""
		expectedPermalink := ""

		// Arrange

		r := ioutil.NopCloser(bytes.NewReader([]byte("")))

		factErrClient := &mocks.MockHttpClient{
			DoFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       r,
				}, errors.New("batata")
			},
		}

		factClient := &BaseClient{
			client: factErrClient,
			url:    "temp"}

		// Act
		response, err := factClient.GetFact()

		// Assert

		assert.NotNil(t, err)

		assert.EqualValues(t,
			expectedText,
			response.Text)

		assert.EqualValues(t,
			expectedId,
			response.ID)

		assert.EqualValues(t,
			expectedSourceUrl,
			response.SourceURL)

		assert.EqualValues(t,
			expectedLanguage,
			response.Language)

		assert.EqualValues(t,
			expectedPermalink,
			response.Permalink)

	})

}

func TestFailedJokeRequest(t *testing.T) {

	t.Run("Error joke request", func(t *testing.T) {

		expectedId := 0
		expectedjoke := ""
		expectedType := ""
		expectedCategories := []string(nil)

		// Arrange
		r := ioutil.NopCloser(bytes.NewReader([]byte("")))

		jokeHttpErrClient := &mocks.MockHttpClient{
			DoFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       r,
				}, errors.New("batata")
			},
		}

		jokeClient := &BaseClient{
			client: jokeHttpErrClient,
			url:    "temp"}

		// Act
		response, err := jokeClient.GetJoke()

		// Assert

		assert.NotNil(t, err)

		assert.EqualValues(t,
			expectedId,
			response.Value.ID)

		assert.EqualValues(t,
			expectedjoke,
			response.Value.Joke)

		assert.EqualValues(t,
			expectedType,
			response.Type)

		assert.EqualValues(t,
			expectedCategories,
			response.Value.Categories)

	})

}

func TestFailedTelegramRequest(t *testing.T) {

	t.Run("Error telegram request", func(t *testing.T) {

		expectedResponseText := ""

		// Arrange
		r := ioutil.NopCloser(bytes.NewReader([]byte(expectedResponseText)))

		telegramHttErrClient := &mocks.MockHttpClient{
			DoFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       r,
				}, errors.New("batata")
			},
		}

		telegramClient := &BaseClient{
			client: telegramHttErrClient,
			url:    "temp"}

		// Act
		response, err := telegramClient.PostResponse(123, "stuff happened")

		// Assert

		assert.NotNil(t, err)

		assert.EqualValues(t,
			expectedResponseText,
			response)

	})

}

func TestSuccessPostTelegram(t *testing.T) {

	t.Run("Successful post to telegram request", func(t *testing.T) {

		expectedResponseText := "{\"ok\": true,\"result\": {\"message_id\": 45,\"from\": {\"id\": 1025326803,\"is_bot\": true,\"first_name\": \"MyDailyFact\",\"username\": \"majoFFper_bot\"},\"chat\": {\"id\": 690639026,\"first_name\": \"Mário\",\"type\": \"private\"},\"date\": 1615076796,\"text\": \"Product Owners never ask Chuck Norris for more features. They ask for mercy.\"}}"

		// Arrange
		r := ioutil.NopCloser(bytes.NewReader([]byte(expectedResponseText)))

		telegramHttSuccessClient := &mocks.MockHttpClient{
			DoFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			},
		}

		telegramClient := &BaseClient{
			client: telegramHttSuccessClient,
			url:    "temp"}

		// Act
		response, err := telegramClient.PostResponse(123, "stuff happened")

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		// Assert

		assert.EqualValues(t,
			expectedResponseText,
			response)

	})

}

func TestSuccessJokeRequest(t *testing.T) {

	t.Run("Successful joke request", func(t *testing.T) {

		expectedId := 479
		expectedjoke := "Chuck Norris does not need to know about class factory pattern. He can instantiate interfaces."
		expectedType := "success"
		expectedCategories := []string([]string{"nerdy"})

		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			if r.Method == http.MethodGet {
				w.Write([]byte("{\"type\": \"success\",\"value\": {\"id\": 479,\"joke\": \"Chuck Norris does not need to know about class factory pattern. He can instantiate interfaces.\",\"categories\": [\"nerdy\"]}}"))
			}
		}))

		jokeClient := &BaseClient{url: ts.URL}

		// Act
		response, err := jokeClient.GetJoke()

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		// Assert

		assert.EqualValues(t,
			expectedId,
			response.Value.ID)

		assert.EqualValues(t,
			expectedjoke,
			response.Value.Joke)

		assert.EqualValues(t,
			expectedType,
			response.Type)

		assert.EqualValues(t,
			expectedCategories,
			response.Value.Categories)

	})

}
func TestSuccessFactRequest(t *testing.T) {

	t.Run("Successful fact request", func(t *testing.T) {

		expectedId := "96221b11-8a37-4495-baf0-134be4feffc1"
		expectedText := "To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P."
		expectedSourceUrl := "http://www.djtech.net/humor/useless_facts.htm"
		expectedLanguage := "en"
		expectedPermalink := "https://uselessfacts.jsph.pl/96221b11-8a37-4495-baf0-134be4feffc1"

		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			if r.Method == http.MethodGet {
				w.Write([]byte("{\"id\": \"96221b11-8a37-4495-baf0-134be4feffc1\", \"text\": \"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.\", \"source\": \"djtech.net\", \"source_url\": \"http://www.djtech.net/humor/useless_facts.htm\", \"language\": \"en\", \"permalink\": \"https://uselessfacts.jsph.pl/96221b11-8a37-4495-baf0-134be4feffc1\"}"))
			}
		}))

		factClient := &BaseClient{url: ts.URL}

		// Act
		response, err := factClient.GetFact()

		if err != nil {
			t.Fatal("Can't run test scenario")
		}

		// Assert

		assert.EqualValues(t,
			expectedText,
			response.Text)

		assert.EqualValues(t,
			expectedId,
			response.ID)

		assert.EqualValues(t,
			expectedSourceUrl,
			response.SourceURL)

		assert.EqualValues(t,
			expectedLanguage,
			response.Language)

		assert.EqualValues(t,
			expectedPermalink,
			response.Permalink)

	})

}
