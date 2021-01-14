package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	// DefaultHTTPGetAddress Default Address
	RandomFactsAddress = "https://uselessfacts.jsph.pl/today.json?language=en"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
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

	update, err := parseTelegramRequest(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Sanitize input
	var sanitizedSeed = sanitize(update.Message.Text)

	generatedFact, err := getRandomFact()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	generateFactText, err := json.Marshal(generatedFact)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// trigger post with response??

	return events.APIGatewayProxyResponse{
		Body:       string(generateFactText),
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

	err = json.Unmarshal([]byte(body), responseGeneratedFact)

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
	return &update, nil
}

// sanitize remove clutter like /start /count or the bot name from the string s passed as input
func sanitize(s string) string {
	if len(s) >= lenStartCommand {
		if s[:lenStartCommand] == startCommand {
			s = s[lenStartCommand:]
		}
	}

	if len(s) >= lenCountCommand {
		if s[:lenCountCommand] == countCommand {
			s = s[lenCountCommand:]
		}
	}
	if len(s) >= lenBotTag {
		if s[:lenBotTag] == botTag {
			s = s[lenBotTag:]
		}
	}
	return s
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
