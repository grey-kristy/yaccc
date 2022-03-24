package convertor

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	API_URL = "https://sandbox-api.coinmarketcap.com/v1"
	API_KEY = "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c"

	MaxBodySize = 10 * 1024 * 1024
)

type apiResp struct {
	Status struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_message"`
	} `json:"Status"`
	Data map[string]json.RawMessage `json:"Data"`
}

type quote struct {
	Price float64 `json:"price"`
}

func prepareReq(amount float64, from_currency string, to_currency string) *http.Request {
	url := API_URL + "/tools/price-conversion"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", API_KEY)

	q := req.URL.Query()
	q.Add("amount", fmt.Sprintf("%f", amount))
	q.Add("symbol", from_currency)
	q.Add("convert", to_currency)
	req.URL.RawQuery = q.Encode()

	return req
}

func getBody(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP error: %s\n", err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, MaxBodySize))
	if err != nil {
		fmt.Printf("HTTP body read error: ", err)
		return []byte{}, err
	}
	return body, nil
}

func processError(errMsg string) {
	errList := strings.Split(errMsg, ":")
	if len(errList) == 2 && (errList[0] == "Invalid value for \"convert\"" ||
		errList[0] == "Invalid value for \"symbol\"") {
		fmt.Printf("Unknown currency symbol: %s\n", errList[1])
	} else {
		fmt.Printf("Error: %s\n", errMsg)
	}
}

func getNextJson(key string, data map[string]json.RawMessage) map[string]json.RawMessage {
	var row map[string]json.RawMessage

	if _, ok := data[key]; !ok {
		fmt.Printf("no key \"%s\" in json: %s", key, data)
		return nil
	}

	if err := json.Unmarshal(data[key], &row); err != nil {
		fmt.Printf("json decode error: ", err)
		return nil
	}

	return row
}

func Convert(amount float64, from_currency string, to_currency string) float64 {
	req := prepareReq(amount, from_currency, to_currency)
	body, err := getBody(req)
	if err != nil {
		return 0
	}
	//fmt.Printf("%s\n", body)

	var data apiResp
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("json decode error: ", err)
		return 0
	}

	if data.Status.ErrorCode != 0 {
		processError(data.Status.ErrorMsg)
		return 0
	}

	var row map[string]json.RawMessage
	if _, ok := data.Data[from_currency]; ok {
		row = getNextJson("quote", getNextJson(from_currency, data.Data))
	} else {
		row = getNextJson("quote", data.Data)
	}

	if row == nil {
		return 0
	}

	var quote quote
	if err := json.Unmarshal(row[to_currency], &quote); err != nil {
		fmt.Printf("json decode error: ", err)
		return 0
	}

	return quote.Price
}
