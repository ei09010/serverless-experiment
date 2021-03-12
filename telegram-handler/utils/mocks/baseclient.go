package mocks

import (
	"my-first-telegram-bot/telegram-handler/restclient"
)

var (
	ReturnGetFact      func() (*restclient.GeneratedFact, error)
	ReturnGetJoke      func() (*restclient.GeneratedJoke, error)
	ReturnPostResponse func(chatId int, text string) (string, error)
)

type MockBaseClient struct {
	ReturnGetFactCallCount      int
	ReturnGetJokeCallCount      int
	ReturnPostResponseCallCount int
}

func (mck *MockBaseClient) GetFact() (*restclient.GeneratedFact, error) {
	mck.ReturnGetFactCallCount++
	return ReturnGetFact()
}

func (mck *MockBaseClient) GetJoke() (*restclient.GeneratedJoke, error) {
	mck.ReturnGetJokeCallCount++
	return ReturnGetJoke()
}
func (mck *MockBaseClient) PostResponse(chatId int, text string) (string, error) {
	mck.ReturnPostResponseCallCount++
	return ReturnPostResponse(chatId, text)
}
