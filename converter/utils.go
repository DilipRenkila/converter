package converter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Request struct {
	Amount   int
	Currency string
}

// used to decode json string obtained from fixer.io
type response struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// decodes json string to a struct.
func decode(r []byte) (x *response, err error) {
	x = new(response)
	err = json.Unmarshal(r, x)
	return
}

// parses amount from query, returns error if amount is not a number
func CheckAmount(s string, dest interface{}) error {
	d, ok := dest.(*int)
	if !ok {
		return fmt.Errorf("bad type for amount: %T", dest)
	}
	n, err := strconv.Atoi(s)
	if err != nil && n > 0 {
		return fmt.Errorf("bad type for amount: %v, give an integer", s)
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
	if len(s)!= 3 {
		return fmt.Errorf("bad type for currency: %s ; should be a three letter string", s)
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return fmt.Errorf("bad type for currency: %s ; should be a three letter string", s)
		}
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

// returns a interface of supported currencies by fixer.io
func KeysString(m map[string]float64) map[string]interface{} {
	keys := make([]string, 0, len(m))
	body := make(map[string]interface{})
	for k := range m {
		keys = append(keys, k)
	}
	body["supported_currencies"] = strings.Join(keys, ", ")

	return body
}

// RunConversion converts the currency by consuming api from fixer.io .
func RunConversion(amount int, currency string) (map[string]interface{},error) {

	body := make(map[string]interface{})
	nestedBody := make(map[string]interface{})

	// creates a http client with a timeout of 5 seconds.
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://api.fixer.io/latest?base=%s", currency)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return body,err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body,err
	}

	values,err := decode(respBody)
	if err != nil {
		return body,err
	}

	body["amount"] = amount
	body["currency"] = currency

	for currencyType, rate := range values.Rates {
		nestedBody[currencyType] = strconv.FormatFloat(round(float64(amount)*rate, 2), 'f', -1, 64)
	}
	body["converted"] = nestedBody
	return body,nil
}