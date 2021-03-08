package restclient

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	Client HTTPClient

	RandomFactsAddress = "https://uselessfacts.jsph.pl/today.json?language=en"

	RandomJokesAddress = "http://api.icndb.com/jokes/random?limitTo=[nerdy]"

	TelegramApi = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_API_TOKEN") + "/sendMessage"

	MyFactClient, MyJokeClient, MyTelegramClient *baseClient
)

type GeneratedFact struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Source    string `json:"source"`
	SourceURL string `json:"source_url"`
	Language  string `json:"language"`
	Permalink string `json:"permalink"`
}

type GeneratedJoke struct {
	Type  string `json:"type"`
	Value struct {
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
		Categories []string `json:"categories"`
	} `json:"value"`
}

type FactClient interface {
	GetFact() (*GeneratedFact, error)
}

type JokeClient interface {
	GetJoke() (*GeneratedJoke, error)
}

type TelegramClient interface {
	PostResponse(chatId int, text string) (string, error)
}

type baseClient struct {
	client http.Client
	url    string
}

func New(url string) *baseClient {
	return &baseClient{
		http.Client{
			Timeout: time.Duration(30) * time.Second,
		},
		url,
	}
}

func (cb *baseClient) GetFact() (*GeneratedFact, error) {

	r, err := get(cb.url)

	responseGeneratedFact := &GeneratedFact{}

	if err != nil {
		return responseGeneratedFact, err
	}

	if r.StatusCode != 200 {
		return responseGeneratedFact, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return responseGeneratedFact, err
	}

	err = json.Unmarshal([]byte(body), &responseGeneratedFact)

	if err != nil {
		return responseGeneratedFact, err
	}

	if len(responseGeneratedFact.Text) == 0 {
		return responseGeneratedFact, err
	}

	return responseGeneratedFact, nil
}

func (cb *baseClient) GetJoke() (*GeneratedJoke, error) {

	r, err := get(cb.url)

	responseGeneratedJoke := &GeneratedJoke{}

	if err != nil {
		return responseGeneratedJoke, err
	}

	if r.StatusCode != 200 {
		return responseGeneratedJoke, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return responseGeneratedJoke, err
	}

	err = json.Unmarshal([]byte(body), &responseGeneratedJoke)

	if err != nil {
		return responseGeneratedJoke, err
	}

	if len(responseGeneratedJoke.Value.Joke) == 0 {
		return responseGeneratedJoke, err
	}

	return responseGeneratedJoke, nil
}

func (cb *baseClient) PostResponse(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)
	response, err := postForm(
		TelegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}

	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	Client = &http.Client{}

	MyJokeClient = New(RandomJokesAddress)

	MyFactClient = New(RandomJokesAddress)

	MyTelegramClient = New(TelegramApi)
}

func get(url string) (*http.Response, error) {

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return Client.Do(request)
}

func post(url, contentType string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return Client.Do(req)
}

func postForm(url string, data url.Values) (resp *http.Response, err error) {

	return post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
