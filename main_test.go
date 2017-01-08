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
	t.Log("checking for missing arguments, bad currency and amount types")
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
		}else{
			t.Logf("Success, Got 400 for %s",url)

		}
	}

}

// TestConversionUnsupportedCurrency checks for conversion from unsupported currency.
func TestConversionUnsupportedCurrency(t *testing.T) {
	t.Log("Checking for conversion of unsupported currency")
	url := fmt.Sprintf("%s/convert?amount=305&currency=USW", server.URL)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/xml")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 400 {
		t.Errorf("Error status code : %d is not expected", res.StatusCode)
	}else{
		t.Logf("Success, Got 400 for %s",url)
	}

}

// TestConversionToJSON checks for Successful listing of converted currencies in JSON.
func TestConversionToJSON(t *testing.T) {
	t.Log("checking for successful listing of converted currencies in JSON")
	url := fmt.Sprintf("%s/convert?amount=305&currency=USD", server.URL)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Error status code : %d is not expected", res.StatusCode)
	}else{
		t.Log("Success, Got 200 for %s",url)
	}


}

// TestConversionToXML checks for Successful listing of converted currencies in XML.
func TestConversionToXML(t *testing.T) {
	t.Log("checking for successful listing of converted currencies in XML ")
	url := fmt.Sprintf("%s/convert?amount=305&currency=USD", server.URL)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/xml")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Error status code : %d is not expected", res.StatusCode)
	}else{
		t.Log("Success, Got 200 for %s",url)
	}

}
