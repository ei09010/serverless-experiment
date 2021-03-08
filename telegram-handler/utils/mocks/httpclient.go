package mocks

import "net/http"

type MockHTTPClient struct {
	DoFuncGET  func(req *http.Request) (*http.Response, error)
	DoFuncPOST func(req *http.Request) (*http.Response, error)
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFuncGET  func(req *http.Request) (*http.Response, error)
	GetDoFuncPOST func(req *http.Request) (*http.Response, error)
)

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {

	if req.Method == http.MethodGet {
		return GetDoFuncGET(req)
	}

	if req.Method == http.MethodPost {
		return GetDoFuncPOST(req)
	}

	return nil, nil
}
