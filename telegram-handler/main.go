package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"my-first-telegram-bot/telegram-handler/restclient"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Configuration struct {
	TelegramApiToken string
}

var (
	botTag = "@CovidCount"

	lenBotTag = len(botTag)

	// Define a few constants and variable to handle different commands
	countCommand = "/givemethecount"

	lenCountCommand = len(countCommand)

	startCommand = "/start"

	lenStartCommand = len(startCommand)

	chatId = 0

	RandomFactsAddress = "https://uselessfacts.jsph.pl/today.json?language=en"

	RandomJokesAddress = "http://api.icndb.com/jokes/random?limitTo=[nerdy]"

	ErrNon200Response = errors.New("Non 200 Response found")

	telegramApi = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_API_TOKEN") + "/sendMessage"
)

type GeneratedFact struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Source    string `json:"source"`
	SourceURL string `json:"source_url"`
	Language  string `json:"language"`
	Permalink string `json:"permalink"`
}

type GenerateJoke struct {
	Type  string `json:"type"`
	Value struct {
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
		Categories []string `json:"categories"`
	} `json:"value"`
}

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("The request has the following body: %s", request.Body)

	update, err := parseTelegramRequest(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Thank you for reaching out, stuff is up and running, but this is telegram bot and this endpoint will eventually cease to exist",
		}, nil
	}

	log.Printf("processed the following text from telegram: %s", update.Message.Text)

	var generatedText []byte

	if strings.Contains(update.Message.Text, "/fact") {

		generatedFact, err := getRandomFact()

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		generatedText, err = json.Marshal(generatedFact.Text)

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
	} else if strings.Contains(update.Message.Text, "/joke") {

		generatedJoke, err := getRandomJoke()

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		generatedText, err = json.Marshal(generatedJoke.Value.Joke)

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
	} else {
		log.Printf("No valid input dected")

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "No valid input dected",
		}, nil
	}

	unquotedStr, err := strconv.Unquote(string(generatedText))

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// send message to telegram through a post
	tempResponse, err := sendTextToTelegramChat(chatId, unquotedStr)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Printf("Got the following response from telegram: %s", tempResponse)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       tempResponse,
	}, nil

}

func getRandomJoke() (GenerateJoke, error) {
	resp, err := restclient.Get(RandomJokesAddress)

	responseGeneratedJoke := GenerateJoke{}

	if err != nil {
		return responseGeneratedJoke, err
	}

	if resp.StatusCode != 200 {
		return responseGeneratedJoke, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

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

func getRandomFact() (GeneratedFact, error) {
	resp, err := restclient.Get(RandomFactsAddress)

	responseGeneratedFact := GeneratedFact{}

	if err != nil {
		return responseGeneratedFact, err
	}

	if resp.StatusCode != 200 {
		return responseGeneratedFact, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

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

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(requestBody string) (*Update, error) {
	var update Update

	if err := json.Unmarshal([]byte(requestBody), &update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}

	chatId = update.Message.Chat.Id

	return &update, nil
}

func sendTextToTelegramChat(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	response, err := restclient.PostForm(
		telegramApi,
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

func main() {
	lambda.Start(handler)
}
