package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"strings"
)

type HTTP struct{}

// tries to post json data to an URL and expects a json return
func (c *HTTP) PostJson(url string, body any, skipSSL bool) ([]byte, error){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
	}
	client := &http.Client{Transport: tr}
	jsonBody, _ := json.Marshal(body)
	res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
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
			return body, nil
		}else{
			return nil, errors.New(fmt.Sprintf("HTTP content-type is not json: %s", contentType))
		}
	}else{
		return nil, errors.New(fmt.Sprintf("HTTP post request failed: %d", res.StatusCode))
	}
}

// tries to get json data from an URL
func (c *HTTP) GetJson(url string, skipSSL bool) ([]byte, error){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(url)
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
			return body, nil
		}else{
			return nil, errors.New(fmt.Sprintf("HTTP content-type is not json: %s", contentType))
		}
	}else{
		return nil, errors.New(fmt.Sprintf("HTTP get request failed: %d", res.StatusCode))
	}
}