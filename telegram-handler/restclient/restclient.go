package restclient

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"my-first-telegram-bot/telegram-handler/dto"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	RandomFactsAddress = "https://uselessfacts.jsph.pl/today.json?language=en"

	RandomJokesAddress = "http://api.icndb.com/jokes/random?limitTo=[nerdy]"

	TelegramApi = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_API_TOKEN") + "/sendMessage"

	MyFactClient FactClient = &BaseClient{
		client: &http.Client{},
		url:    RandomFactsAddress,
	}

	MyJokeClient JokeClient = &BaseClient{
		client: &http.Client{},
		url:    RandomJokesAddress,
	}

	MyTelegramClient TelegramClient = &BaseClient{
		client: &http.Client{},
		url:    TelegramApi,
	}
)

type FactClient interface {
	GetFact() (*dto.GeneratedFact, error)
}

type JokeClient interface {
	GetJoke() (*dto.GeneratedJoke, error)
}

type TelegramClient interface {
	PostResponse(chatId int, content string) (string, error)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type BaseClient struct {
	client HttpClient
	url    string
}

func (cb *BaseClient) GetFact() (*dto.GeneratedFact, error) {

	r, err := get(cb)

	factToReturn := &dto.GeneratedFact{}

	if err != nil {
		return factToReturn, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return factToReturn, err
	}

	err = json.Unmarshal([]byte(body), &factToReturn)

	if err != nil {
		return factToReturn, err
	}

	if len(factToReturn.Text) == 0 {
		return factToReturn, err
	}

	return factToReturn, nil
}

func (cb *BaseClient) GetJoke() (*dto.GeneratedJoke, error) {

	r, err := get(cb)

	jokeToReturn := &dto.GeneratedJoke{}

	if err != nil {
		return jokeToReturn, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return jokeToReturn, err
	}

	err = json.Unmarshal([]byte(body), &jokeToReturn)

	if err != nil {
		return jokeToReturn, err
	}

	if len(jokeToReturn.Value.Joke) == 0 {
		return jokeToReturn, err
	}

	return jokeToReturn, nil
}

func (cb *BaseClient) PostResponse(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)

	response, err := postForm(
		cb,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)

	if errRead != nil {

		return "", err
	}

	bodyString := string(bodyBytes)

	return bodyString, nil
}

func get(bc *BaseClient) (*http.Response, error) {

	request, err := http.NewRequest(http.MethodGet, bc.url, nil)

	if err != nil {
		return nil, err
	}

	return bc.client.Do(request)
}

func post(cb *BaseClient, contentType string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPost, cb.url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return cb.client.Do(req)
}

func postForm(cb *BaseClient, data url.Values) (resp *http.Response, err error) {

	return post(cb, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

}
