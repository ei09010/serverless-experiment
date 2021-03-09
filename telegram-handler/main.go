package main

import (
	"encoding/json"
	"errors"
	"log"
	"my-first-telegram-bot/telegram-handler/restclient"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	chatId = 0

	ErrNon200Response = errors.New("Non 200 Response found")
)

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

		generatedFact, err := restclient.MyFactClient.GetFact()

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		generatedText, err = json.Marshal(generatedFact.Text)

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
	} else if strings.Contains(update.Message.Text, "/joke") {

		generatedJoke, err := restclient.MyFactClient.GetJoke()

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
	tempResponse, err := restclient.MyTelegramClient.PostResponse(chatId, unquotedStr)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Printf("Got the following response from telegram: %s", tempResponse)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       tempResponse,
	}, nil

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

func main() {
	lambda.Start(handler)
}
