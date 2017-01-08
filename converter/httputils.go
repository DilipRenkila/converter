package converter

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Diggernaut/mxj"
)

// outputWriter writes output to http response in specified format.
func outputWriter(w http.ResponseWriter, body map[string]interface{}, contentType string, status int) (http.ResponseWriter, error) {
	// checks for the output version, either in XML or JSON
	re, err := regexp.Compile(`xml`)
	if re.MatchString(contentType) == true {
		w, err = xmlWriter(w, body, status)
		return w, err
	}
	w, err = jsonWriter(w, body, status)
	return w, err
}

// jsonWriter writes the output to the http response stream as json with standard json encoding.
func jsonWriter(w http.ResponseWriter, body map[string]interface{}, status int) (http.ResponseWriter, error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	return w, err
}

// xmlWriter writes the output to the http response stream as xml with standard xml encoding.
func xmlWriter(w http.ResponseWriter, body map[string]interface{}, code int) (http.ResponseWriter, error) {

	data, err := mxj.AnyXmlIndentByte(body, "", " ")
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	w.Write(data)
	return w, err

}
