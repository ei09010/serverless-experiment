package mocks

import (
	"my-first-telegram-bot/telegram-handler/restclient"
)

var (
	ReturnGetFact      func() (*restclient.GeneratedFact, error)
	ReturnGetJoke      func() (*restclient.GeneratedJoke, error)
	ReturnPostResponse func(chatId int, text string) (string, error)
)

type MockBaseClient struct{}

func (mck *MockBaseClient) GetFact() (*restclient.GeneratedFact, error) {
	return ReturnGetFact()
}

func (mck *MockBaseClient) GetJoke() (*restclient.GeneratedJoke, error) {
	return ReturnGetJoke()
}
func (mck *MockBaseClient) PostResponse(chatId int, text string) (string, error) {
	return ReturnPostResponse(chatId, text)
}
