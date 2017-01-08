package converter

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// ConvertCurrency is a HTTP handler to construct an HTTP response for route /convert.
func ConvertCurrency(w http.ResponseWriter, r *http.Request) {

	// body a map of strings to arbitrary data types used to construct the response.
	body := make(map[string]interface{})

	// contentType stores the first value associated with the given key "Accept" in http request header.
	contentType := r.Header.Get("Accept")

	// Parse the remote client information.
	addr, proto, err := RestClientIP(r)
	if err != nil {
		log.Printf("Error in parsing Client's IP: %s", err)
		body["reason"] = "Bad request type"
		w, err = outputWriter(w, body, contentType, http.StatusBadRequest)
		if err != nil {
			log.Printf("Error in encoding: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Request came from: %s by %s", addr, proto)

	// ParseForm parses the raw Request from the URL, returns err if it fails to parse.
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		body["reason"] = "Bad request type"
		w, err = outputWriter(w, body, contentType, http.StatusBadRequest)
		if err != nil {
			log.Printf("Error in encoding: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// checks whether api client has given amount and returns StatusBadRequest if they aren't given.
	if r.Form.Get("amount") == "" {
		body["reason"] = "missing arguments amount"
		w, err = outputWriter(w, body, contentType, http.StatusBadRequest)
		if err != nil {
			log.Printf("Error in encoding: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// checks whether api client has given currency and returns StatusBadRequest if they aren't given.
	if r.Form.Get("currency") == "" {
		body["reason"] = "missing arguments currency"
		w, err = outputWriter(w, body, contentType, http.StatusBadRequest)
		if err != nil {
			log.Printf("Error in encoding: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	money := &Request{}

	// check errors in amount given by api client.
	if err := CheckAmount(r.Form.Get("amount"), &money.Amount); err != nil {
		body["reason"] = fmt.Sprintf("%v", err)
		w, err = outputWriter(w, body, contentType, http.StatusBadRequest)
		if err != nil {
			log.Printf("Error in encoding: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// check errors in currency type given by api client.
	if err := CheckCurrency(r.Form.Get("currency"), &money.Currency); err != nil {
		body["reason"] = fmt.Sprintf("%v", err)
		w, err = outputWriter(w, body, contentType, http.StatusBadRequest)
		if err != nil {
			log.Printf("Error in encoding: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	body, err = RunConversion(money.Amount, money.Currency)
	if err != nil {

		log.Printf("%v", err)
		body["reason"] = "Internal Server Error"
		w, err = outputWriter(w, body, contentType, http.StatusInternalServerError)
		if err != nil {
			log.Printf("Error in encoding: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w, err = outputWriter(w, body, contentType, http.StatusOK)
	if err != nil {
		log.Printf("Error in encoding: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
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
