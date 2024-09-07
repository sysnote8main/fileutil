package httputil

import "net/http"

var (
	client = &http.Client{}
)

func GetDefaultHttpClient() *http.Client {
	return client
}

func NewHttpClient() *http.Client {
	return &http.Client{}
}

func NewGetRequest(url string, headers map[string][]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return req, nil
}
