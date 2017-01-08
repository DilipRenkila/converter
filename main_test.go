package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server   *httptest.Server
	userURLs []string
)

func init() {
	//Creating a new http server with the user handlers.
	server = httptest.NewServer(NewRouter())

}

// TestArguments checks for missing arguments, bad currency and amount types.
func TestArguments(t *testing.T) {

	userURLs = []string{fmt.Sprintf("%s/convert", server.URL),
		fmt.Sprintf("%s/convert?", server.URL),
		fmt.Sprintf("%s/convert?amount=abcd", server.URL),
		fmt.Sprintf("%s/convert?amount=-59", server.URL),
		fmt.Sprintf("%s/convert?currency=njkkkk", server.URL),
		fmt.Sprintf("%s/convert?currency=se&amount=24", server.URL),
	}

	for _, url := range userURLs {
		request, err := http.NewRequest("GET", url, nil)
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != 400 {
			t.Errorf("Bad Request Type expected: %d", res.StatusCode)
		}
	}

}

// TestConversionUnsupportedCurrency checks for conversion from unsupported currency.
func TestConversionUnsupportedCurrency(t *testing.T) {

	url := fmt.Sprintf("%s/convert?amount=305&currency=USW", server.URL)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/xml")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 400 {
		t.Errorf("Error status code : %d is not expected", res.StatusCode)
	}

}

// TestConversionToJSON checks for Successful listing of converted currencies in XML.
func TestConversionToJSON(t *testing.T) {

	url := fmt.Sprintf("%s/convert?amount=305&currency=USD", server.URL)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Error status code : %d is not expected", res.StatusCode)
	}

}

// TestConversionToXML checks for Successful listing of converted currencies in XML.
func TestConversionToXML(t *testing.T) {

	url := fmt.Sprintf("%s/convert?amount=305&currency=USD", server.URL)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/xml")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Error status code : %d is not expected", res.StatusCode)
	}

}
