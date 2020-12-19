package internal

import (
	"encoding/json"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"time"
	"crypto/hmac"
	"crypto/sha512"
	"io"
	"math/rand"
	"strings"
)

var funds = Funds{}


func nonce() string {
	rand.Seed(time.Now().Unix())
	var output strings.Builder
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	length := 20
	for i := 0; i < length; i++ {
		output.WriteString(string(charSet[rand.Intn(len(charSet))]))
	}
	return output.String()
}

func sign(message string, secret string) string {
	hash := hmac.New(sha512.New, []byte(secret))
	io.WriteString(hash, message)
	return hex.EncodeToString(hash.Sum(nil))
}

func addHeaders(url string, key string, secret string, req *http.Request) {
	req.Header.Set("User-Agent", "PyRock v1")
	req.Header.Set("content-type", "application/json")
	if len(key) > 0 {
		nonce := nonce()
		message := nonce + url
		req.Header.Set("X-TRT-KEY", key)
		req.Header.Set("X-TRT-SIGN", sign(message, secret))
		req.Header.Set("X-TRT-NONCE", nonce)
	}
}

func get(url string, key string, secret string) ([]byte, error) {
	WaitDueRateLimit(RATE_LIMIT)
	client := http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	addHeaders(url, key, secret, req)
	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	return body, nil
}

func GetFunds(cred Credential) (Funds, error) {
	if len(funds.Funds) <1 {
		body, err := get(API_BASEURL + "funds", cred.RoKey, cred.RoSecret)
		if err != nil {
			return funds, err
		}
		err = json.Unmarshal(body, &funds)
		if err != nil {
			return funds, err
		}
	}
	return funds, nil
}

func GetTicker(fund string, cred Credential) (Ticker, error) {
	ticker := Ticker{}
	body, err := get(API_BASEURL + "funds/" + fund + "/ticker", cred.RoKey, cred.RoSecret)
	if err != nil {
		return ticker, err
	}
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		return ticker, err
	}
	return ticker, nil
}

func GetTickers(cred Credential) (Tickers, error) {
	tickers := Tickers{}
	body, err := get(API_BASEURL + "funds/tickers", cred.RoKey, cred.RoSecret)
	if err != nil {
		return tickers, err
	}
	err = json.Unmarshal(body, &tickers)
	if err != nil {
		return tickers, err
	}
	return tickers, nil
}