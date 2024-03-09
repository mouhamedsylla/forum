package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"syscall/js"
	"time"
)

var (
	Forum_Categories = map[string]bool{
		"Informatique":         true,
		"Software Engineering": true,
		"Education":            true,
		"kitchen":              true,
	}
)

// The ParseComponent function in Go parses a template string and executes it with the provided data,
// returning the result as a string.
func ParseComponent(component string, p interface{}) string {
	var resultBuffer bytes.Buffer
	tpml, err1 := template.New("template").Parse(component)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	err := tpml.Execute(&resultBuffer, p)
	if err != nil {
		panic(err)
	}
	return resultBuffer.String()
}

// The SendData function sends JSON data to a specified endpoint asynchronously and calls a callback
// function with the response or error.
func SendData(data interface{}, endpoint string, callback func(*http.Response, error)) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		callback(nil, err)
	}
	go func() {
		url := "http://localhost:8080/" + endpoint
		resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(jsonData))
		callback(resp, err)
	}()
}

// The `GetData` function asynchronously fetches data from a specified endpoint URL and invokes a
// callback function with the response and any errors.
func GetData(endpoint string, callback func(*http.Response, error)) {
	go func() {
		url := "http://localhost:8080/" + endpoint
		resp, err := http.Get(url)
		callback(resp, err)
	}()
}

// The function `DecodeResponse` reads and decodes the JSON response body from an HTTP response into a
// target interface.
func DecodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

func IsValideCategories(categories []string) bool {
	for _, c := range categories {
		_, ok := Forum_Categories[c]
		if !ok {
			return false
		}
	}
	return true
}

func NotNullEntry(data ...string) bool {
	for _, str := range data {
		if len(strings.TrimSpace(str)) == 0 {
			return false
		}
	}
	return true
}

func CheckFinishSession(TimeOut time.Time) {
	b := TimeOut.Before(time.Now())
	if b {
		js.Global().Get("window").Get("location").Set("href", "/connect")
	}
}
