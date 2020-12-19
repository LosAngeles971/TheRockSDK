package internal

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"bytes"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

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

func Get(url string, key string, secret string) ([]byte, error) {
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
	if os.Getenv("THEROCK_DEBUG") == "true" {
		log.Println("Got: " + url)
		log.Println(string(body))
	}
	return body, nil
}

func Delete(url string, key string, secret string) ([]byte, error) {
	WaitDueRateLimit(RATE_LIMIT)
	client := http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
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
	if os.Getenv("THEROCK_DEBUG") == "true" {
		log.Println("Got: " + url)
		log.Println(string(body))
	}
	return body, nil
}

func Post(url string, data []byte, key string, secret string) ([]byte, error) {
	WaitDueRateLimit(RATE_LIMIT)
	client := http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
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
	if os.Getenv("THEROCK_DEBUG") == "true" {
		log.Println("Got: " + url)
		log.Println(string(body))
	}
	return body, nil
}