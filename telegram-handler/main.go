package main

import (
	"encoding/json"
	"errors"
	"log"
	"my-first-telegram-bot/telegram-handler/dto"
	"my-first-telegram-bot/telegram-handler/restclient"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	chatId                      = 0
	TELEGRAM_FACT_REQUEST_TOKEN = "/fact"
	TELEGRAM_JOKE_REQUEST_TOKEN = "/joke"

	ErrNon200Response        = errors.New("Non 200 Response found")
	ErrorHttpRequest         = "Error executing http request"
	InformalInvalidResponse  = "Thank you for reaching out, stuff is up and running, but this is a telegram bot and this endpoint will eventually vanish"
	InvalidInputFromTelegram = "No valid input from telegram request detected"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("The request has the following body: %s", request.Body)

	update, err := parseTelegramRequest(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       InformalInvalidResponse,
		}, nil
	}

	var generatedText []byte

	if strings.Contains(update.Message.Text, TELEGRAM_FACT_REQUEST_TOKEN) {

		generatedFact, err := restclient.MyFactClient.GetFact()

		if err != nil {

			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       ErrorHttpRequest,
			}, err

		}

		generatedText, err = json.Marshal(generatedFact.Text)

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

	} else if strings.Contains(update.Message.Text, TELEGRAM_JOKE_REQUEST_TOKEN) {

		generatedJoke, err := restclient.MyJokeClient.GetJoke()

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       ErrorHttpRequest,
			}, err
		}

		generatedText, err = json.Marshal(generatedJoke.Value.Joke)

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

	} else {
		log.Printf(InvalidInputFromTelegram)

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       InvalidInputFromTelegram,
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

func parseTelegramRequest(requestBody string) (*dto.Update, error) {
	var update dto.Update

	if err := json.Unmarshal([]byte(requestBody), &update); err != nil {

		return nil, err
	}

	chatId = update.Message.Chat.Id

	return &update, nil
}

func main() {
	lambda.Start(handler)
}
