package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	botTag = "@CovidCount"

	lenBotTag = len(botTag)

	// Define a few constants and variable to handle different commands
	countCommand = "/givemethecount"

	lenCountCommand = len(countCommand)

	startCommand = "/start"

	lenStartCommand = len(startCommand)

	chatId = 0

	// DefaultHTTPGetAddress Default Address
	RandomFactsAddress = "https://uselessfacts.jsph.pl/today.json?language=en"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")

	telegramApi = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
)

type GeneratedFact struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Source    string `json:"source"`
	SourceURL string `json:"source_url"`
	Language  string `json:"language"`
	Permalink string `json:"permalink"`
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

	log.Println("The request has the following body: %V", request.Body)

	if request.HTTPMethod != http.MethodPost {
		log.Println("Received irrelevant request")

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Thank you for reaching out, stuff is up and running, but this is telegram bot and this endpoint will eventually cease to exist",
		}, nil
	}
	update, err := parseTelegramRequest(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Println(update.Message.Text)

	generatedFact, err := getRandomFact()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	generateFactText, err := json.Marshal(generatedFact.Text)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	unquotedStr, err := strconv.Unquote(string(generateFactText))

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// send message to telegram through a post
	tempResponse, err := sendTextToTelegramChat(chatId, unquotedStr)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Println(tempResponse)

	// returns response to caller
	return events.APIGatewayProxyResponse{
		Body:       unquotedStr,
		StatusCode: 200,
	}, nil

}

func getRandomFact() (GeneratedFact, error) {
	resp, err := http.Get(RandomFactsAddress)

	responseGeneratedFact := GeneratedFact{}

	if err != nil {
		return responseGeneratedFact, err
	}

	if resp.StatusCode != 200 {
		return responseGeneratedFact, err
	}

	if err != nil {
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
	response, err := http.PostForm(
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
