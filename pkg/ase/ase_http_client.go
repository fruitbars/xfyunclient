package ase

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	"xfyunclient/pkg/utils"
)

type ASEHttpClient struct {
	*ASEClientBase
}

// NewASEClient initializes a new ASEClient instance
func NewASEClient(serverURL, appid, apikey, apisecret, httpProto, aseAlgorithm string) *ASEHttpClient {
	baseClient := NewASEClientBase(serverURL, appid, apikey, apisecret, httpProto, aseAlgorithm)
	return &ASEHttpClient{ASEClientBase: baseClient}
}

// createASERequest creates an HTTP request with the necessary headers and body
func (c *ASEHttpClient) createASERequest(reqJsonByte []byte) (*http.Request, error) {
	reqBody := bytes.NewReader(reqJsonByte)
	callURL, headers, err := c.getAuthServerURL("POST")
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", callURL, reqBody)
	if err != nil {
		return nil, err
	}

	request.Header = *headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json,version=1.0")

	//log.Println(request.Header)

	return request, nil
}

func (c *ASEHttpClient) CallASEAPIJson(req interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.CallASEAPI(jsonData)
}

// CallASEAPI sends the request to the ASE API and returns the response body
func (c *ASEHttpClient) CallASEAPI(reqJsonByte []byte) ([]byte, error) {
	request, err := c.createASERequest(reqJsonByte)
	if err != nil {
		return nil, err
	}

	utils.PrintHttpRequestHeader(request)

	client := &http.Client{Timeout: c.timeOut}
	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//log.Println(string(body))

	elapsed := time.Since(start)
	log.Printf("Function took %v s to execute.", elapsed.Seconds())

	return body, nil
}
