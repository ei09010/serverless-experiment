package main

import (
	"encoding/json"
	"errors"
	"log"
	dto "my-first-telegram-bot/telegram-handler/Dto"
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("The request has the following body: %s", request.Body)

	update, err := parseTelegramRequest(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Thank you for reaching out, stuff is up and running, but this is a telegram bot and this endpoint will eventually vanish",
		}, nil
	}

	log.Printf("processed the following text from telegram: %s", update.Message.Text)

	var generatedText []byte

	if strings.Contains(update.Message.Text, "/fact") {

		generatedFact, err := restclient.MyFactClient.GetFact()

		if err != nil {

			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error calling fact generation api",
			}, err

		}

		generatedText, err = json.Marshal(generatedFact.Text)

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

	} else if strings.Contains(update.Message.Text, "/joke") {

		generatedJoke, err := restclient.MyJokeClient.GetJoke()

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error calling joke generation api",
			}, err
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

	tempResponse, err := restclient.MyTelegramClient.PostResponse(chatId, unquotedStr)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       tempResponse,
		}, err
	}

	log.Printf("Got the following response from telegram: %s", tempResponse)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       tempResponse,
	}, nil

}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(requestBody string) (*dto.Update, error) {
	var update dto.Update

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
