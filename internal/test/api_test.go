package test

import (
	"fmt"
	"net/http"
	"testing"
)

func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func TestAPI(t *testing.T) {
	res, err := Get("http://localhost:4200/v1/users")

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(res)
}
