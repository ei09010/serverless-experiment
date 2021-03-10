package restclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestClient(t *testing.T) {

	t.Run("Successful fact request", func(t *testing.T) {

		expectedId := "96221b11-8a37-4495-baf0-134be4feffc1"
		expectedText := "To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P."
		expectedSourceUrl := "http://www.djtech.net/humor/useless_facts.htm"
		expectedLanguage := "en"
		expectedPermalink := "https://uselessfacts.jsph.pl/96221b11-8a37-4495-baf0-134be4feffc1"

		// Arrange
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			if r.Method == http.MethodGet {
				w.Write([]byte("{\"id\": \"96221b11-8a37-4495-baf0-134be4feffc1\", \"text\": \"To Ensure Promptness, one is expected to pay beyond the value of service – hence the later abbreviation: T.I.P.\", \"source\": \"djtech.net\", \"source_url\": \"http://www.djtech.net/humor/useless_facts.htm\", \"language\": \"en\", \"permalink\": \"https://uselessfacts.jsph.pl/96221b11-8a37-4495-baf0-134be4feffc1\"}"))
			}
		}))

		factClient := &BaseClient{url: ts.URL}

		// Act
		response, err := factClient.GetFact()
		if err != nil {
			t.Fatal("Everything should be ok")
		}

		// Assert

		assert.EqualValues(t,
			expectedText,
			response.Text)

		assert.EqualValues(t,
			expectedId,
			response.ID)

		assert.EqualValues(t,
			expectedSourceUrl,
			response.SourceURL)

		assert.EqualValues(t,
			expectedLanguage,
			response.Language)

		assert.EqualValues(t,
			expectedPermalink,
			response.Permalink)

	})

}
