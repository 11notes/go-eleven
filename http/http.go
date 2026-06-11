package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type HTTP struct{}

// tries to post json data to an URL and expects a json return
func (c *HTTP) PostJson(url string, body any) (any, error){
	jsonBody, _ := json.Marshal(body)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 100 && res.StatusCode <= 399 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		contentType := strings.Join(res.Header["Content-Type"], "")
		if strings.HasPrefix(contentType, "application/json") {
			var result any
			json.Unmarshal([]byte(body), &result)
			return result, nil
		}else{
			return nil, errors.New(fmt.Sprintf("HTTP content-type is not json: %s", contentType))
		}
	}else{
		return nil, errors.New(fmt.Sprintf("HTTP post request failed: %d", res.StatusCode))
	}
}

// tries to get json data from an URL
func (c *HTTP) GetJson(url string) (any, error){
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 100 && res.StatusCode <= 399 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		contentType := strings.Join(res.Header["Content-Type"], "")
		if strings.HasPrefix(contentType, "application/json") {
			var result any
			json.Unmarshal([]byte(body), &result)
			return result, nil
		}else{
			return nil, errors.New(fmt.Sprintf("HTTP content-type is not json: %s", contentType))
		}
	}else{
		return nil, errors.New(fmt.Sprintf("HTTP get request failed: %d", res.StatusCode))
	}
}