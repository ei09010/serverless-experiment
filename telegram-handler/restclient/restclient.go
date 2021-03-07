package restclient

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	Client HTTPClient
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	Client = &http.Client{}
}

func Get(url string) (*http.Response, error) {

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

func PostForm(url string, data url.Values) (resp *http.Response, err error) {

	return post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
