package converter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"log"
)

// Request is composed of a amount and currency used to store from http request.
type Request struct {
	Amount   int
	Currency string
}

// Response stores decoded json string obtained from fixer.io.
type Response struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

//decodes json string to a struct of type Response.
func decode(r []byte) (x *Response, err error) {
	x = new(Response)
	err = json.Unmarshal(r, x)
	return
}

// CheckAmount parses amount from query, returns error if amount is not a number
func CheckAmount(s string, dest interface{}) error {
	d, ok := dest.(*int)
	if !ok {
		return fmt.Errorf("bad type for amount: %T", dest)
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return fmt.Errorf("bad type for amount: %v, give a positive integer", s)
	}
	*d = n
	return nil
}

// CheckCurrency parses currency from query and returns error if currency is not a three letter string.
func CheckCurrency(s string, dest interface{}) error {
	d, ok := dest.(*string)
	if !ok {
		return fmt.Errorf("bad type for currency: %T", dest)
	}
	if len(s) != 3 {
		return fmt.Errorf("bad type for currency: %s ; should be a three letter string", s)
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return fmt.Errorf("bad type for currency: %s ; should be a three letter string", s)
		}
	}

	// checks whether currency given by api client is supported or not.
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://api.fixer.io/latest?base=%s", strings.ToUpper(s))
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	Resp, err := client.Do(req)
	if err != nil {

		log.Printf("Error: %v",err)
		return fmt.Println("Internal Server Error")
	}

	defer Resp.Body.Close()
	if Resp.StatusCode != 200 {
		url := fmt.Sprintf("http://api.fixer.io/latest")
		req, err := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)

		if err != nil {
			log.Printf("Error: %v",err)
			return fmt.Println("Internal Server Error")
		}

		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error: %v",err)
			return fmt.Println("Internal Server Error")
		}
		x,err := decode(resp_body)
		if err != nil {
			log.Printf("Error: %v",err)
			return fmt.Println("Internal Server Error")
		}
		return  fmt.Errorf("Requested Currency type : %s  is not supported. Select from %s", strings.ToUpper(s),KeysToString(x.Rates))
	}

	*d = strings.ToUpper(s)
	return nil
}

// rounds the floating number upto a given precision
func round(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(ret(num*output)) / output

}

func ret(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// KeysToString returns a interface of supported currencies by fixer.io
func KeysToString(m map[string]float64) map[string]interface{} {
	keys := make([]string, 0, len(m))
	body := make(map[string]interface{})
	for k := range m {
		keys = append(keys, k)
	}
	body["supported_currencies"] = strings.Join(keys, ", ")

	return body
}

// RunConversion converts the currency by consuming api from fixer.io .
func RunConversion(amount int, currency string) (map[string]interface{}, error) {

	body := make(map[string]interface{})
	nestedBody := make(map[string]interface{})

	// creates a http client with a timeout of 5 seconds.
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://api.fixer.io/latest?base=%s", currency)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	Resp, err := client.Do(req)
	if err != nil {
		return body, err
	}

	defer Resp.Body.Close()
	RespBody, err := ioutil.ReadAll(Resp.Body)
	if err != nil {
		return body, err
	}

	values, err := decode(RespBody)
	if err != nil {
		return body, err
	}

	body["amount"] = amount
	body["currency"] = currency

	for currencyType, rate := range values.Rates {
		nestedBody[currencyType] = strconv.FormatFloat(round(float64(amount)*rate, 2), 'f', -1, 64)
	}
	body["converted"] = nestedBody
	return body, nil
}

// RestClientIP parses IP configuration of client, returns IP address of the client,protocol (IPv4 or IPv6) and a error if any.
func RestClientIP(r *http.Request) (string, string, error) {
	var address string
	var protocol string

	viaProxy := r.Header.Get("X-Forwarded-For")
	if viaProxy != "" {
		address = viaProxy
	} else {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			address = host
		} else {
			address = r.RemoteAddr
		}
	}

	ip := net.ParseIP(address)
	if ip == nil {
		return "", "", fmt.Errorf("Invalid address: %s", address)
	}

	if ip.To4() == nil {
		protocol = "IPv6"
	} else {
		protocol = "IPv4"
	}

	return address, protocol, nil
}
