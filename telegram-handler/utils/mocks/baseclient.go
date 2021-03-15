package mocks

import (
	dto "my-first-telegram-bot/telegram-handler/Dto"
	"net/http"
)

var (
	ReturnGetFact      func() (*dto.GeneratedFact, error)
	ReturnGetJoke      func() (*dto.GeneratedJoke, error)
	ReturnPostResponse func(chatId int, text string) (string, error)
)

type MockBaseClient struct {
	ReturnGetFactCallCount      int
	ReturnGetJokeCallCount      int
	ReturnPostResponseCallCount int
}

func (mck *MockBaseClient) GetFact() (*dto.GeneratedFact, error) {
	mck.ReturnGetFactCallCount++
	return ReturnGetFact()
}

func (mck *MockBaseClient) GetJoke() (*dto.GeneratedJoke, error) {
	mck.ReturnGetJokeCallCount++
	return ReturnGetJoke()
}
func (mck *MockBaseClient) PostResponse(chatId int, text string) (string, error) {
	mck.ReturnPostResponseCallCount++
	return ReturnPostResponse(chatId, text)
}

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (mckHt *MockHttpClient) Do(req *http.Request) (*http.Response, error) {

	return mckHt.DoFunc(req)
}
