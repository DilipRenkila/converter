## Currency Converter

This repository provides an http api written in Golang that takes 
the price in one currency and converts it to other currencies.

This api provides a single endpoint `/convert?amount=<amount>&currency=<currency>` 
that wraps up a call to fixer.io  to return a JSON or XML array of amounts converted to other currencies for a given ammount and currency.

### Project Structure

```
├── converter
│   ├── handlers.go          - provides http handler functions.
│   ├── httputils.go         - provides XML and JSON wrapper functions. 
│   └── utils.go             - provides internal functions for calling the fixer.io API.   
├── main.go                  - provides the server and handler initalization.
├── main_test.go
├── README.md
└── shipwallet.iml

```

### Getting Started

####Install Prerequisites
- `go get github.com/diliprenkila/converter/converter`
- `go get github.com/gorilla/mux`
- `go get github.com/Diggernaut/mxj`
- `go run main.go`

#### Running the tests

- `go test main.go`

#### Example (json, default)

      Request:

      curl -i 'http://localhost:8080/convert?amount=200&currency=SEK'

      Response:
        
        HTTP/1.1 200 OK
        Content-Type: application/json
        Date: Wed, 04 Jan 2017 12:42:54 GMT
        Content-Length: 496

        {"amount":200,"converted":{"AUD":"30.18","BGN":"41.02","BRL":"71.14","CAD":"29.29","CHF":"22.45","CNY":"151.62","CZK":"566.76","DKK":"155.94","EUR":"20.98","GBP":"17.74","HKD":"168.97","HRK":"158.69","HUF":"6480.2","IDR":"293620","ILS":"84.23","INR":"1489.3","JPY":"2574.8","KRW":"26292","MXN":"450.86","MYR":"97.86","NOK":"188.81","NZD":"31.5","PHP":"1085.82","PLN":"92.16","RON":"94.87","RUB":"1319.62","SGD":"31.59","THB":"782.12","TRY":"78.16","USD":"21.78","ZAR":"299.62"},"currency":"SEK"}

### Example (xml)

      Request:

      curl -H "Accept: application/xml, */*" -i 'http://localhost:8080/convert?amount=200&currency=SEK'

      Response:

        HTTP/1.1 200 OK
        Content-Type: application/xml
        Date: Wed, 04 Jan 2017 12:45:18 GMT
        Content-Length: 692

        <doc>
        <amount>200</amount>
        <converted>
        <AUD>30.18</AUD>
        <BGN>41.02</BGN>
        <BRL>71.14</BRL>
        <CAD>29.29</CAD>
        <CHF>22.45</CHF>
        <CNY>151.62</CNY>
        <CZK>566.76</CZK>
        <DKK>155.94</DKK>
        <EUR>20.98</EUR>
        <GBP>17.74</GBP>
        <HKD>168.97</HKD>
        <HRK>158.69</HRK>
        <HUF>6480.2</HUF>
        <IDR>293620</IDR>
        <ILS>84.23</ILS>
        <INR>1489.3</INR>
        <JPY>2574.8</JPY>
        <KRW>26292</KRW>
        <MXN>450.86</MXN>
        <MYR>97.86</MYR>
        <NOK>188.81</NOK>
        <NZD>31.5</NZD>
        <PHP>1085.82</PHP>
        <PLN>92.16</PLN>
        <RON>94.87</RON>
        <RUB>1319.62</RUB>
        <SGD>31.59</SGD>
        <THB>782.12</THB>
        <TRY>78.16</TRY>
        <USD>21.78</USD>
        <ZAR>299.62</ZAR>
        </converted>
        <currency>SEK</currency>

### Example (json, error handling)
        Request:
        
        curl -i 'http://localhost:8080/convert?amount=200&currency=SE'
        
        Response:
        
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Wed, 04 Jan 2017 12:48:51 GMT
        Content-Length: 68

        {"reason":"Bad Currency type:SE ; should be a three letter string"}
        
        
        Request:
        
        curl -i 'http://localhost:8080/convert?amount=200&currency=SEi'
        
        Response:
        
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Wed, 04 Jan 2017 12:51:04 GMT
        Content-Length: 269

        {"reason":"SEI currency type is not supported and supported currencies are listed below","supported_currencies":"HKD, USD, ZAR, HUF, ILS, MYR, CNY, KRW, SEK, CAD, DKK, NOK, BRL, HRK, PHP, CZK, IDR, JPY, RUB, SGD, AUD, BGN, CHF, THB, TRY, GBP, NZD, PLN, INR, MXN, RON"}

        Request:
        
        curl -i 'http://localhost:8080/convert?amount=200'
        
        Response:
        
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Wed, 04 Jan 2017 12:52:10 GMT
        Content-Length: 66

        {"reason":"Bad Currency type: ; should be a three letter string"}


My solution handles errors effectively by either returning 400 or 500 status codes.

If there is a connection problem like your server is unable to reach http://fixer.io/ or it exceeds timeout of 5 seconds, it responses with Internal error status.
 