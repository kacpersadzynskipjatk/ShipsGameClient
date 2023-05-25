package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ServerAddress   = "https://go-pjatk-server.fly.dev/api"
	GameEndpoint    = "/game"
	BoardEndpoint   = "/game/board"
	FireEndpoint    = "/game/fire"
	DescEndpoint    = "/game/desc"
	ListEndpoint    = "/game/list"
	AbandonEndpoint = "/game/abandon"
	RefreshEndpoint = "/game/refresh"
)

type Client struct {
	HttpClient http.Client
	Token      string
}

func NewGameClient(c *http.Client) *Client {
	gc := &Client{
		HttpClient: *c,
	}
	return gc
}

func (gc *Client) doRequest(requestType, address, token string, payload interface{}) (*http.Response, error) {
	var req *http.Request
	var err error

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(requestType, address, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(requestType, address, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-auth-token", token)
	resp, err := gc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (gc *Client) SendRequest(methodType, endpoint, token string, request interface{}, response interface{}) (interface{}, error) {
	for n := 1; n <= 10; n++ {
		resp, err := gc.doRequest(methodType, ServerAddress+endpoint, token, &request)

		if err != nil {
			return response, err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			if err != io.EOF {
				return response, err
			}
		}
		gc.Token = resp.Header.Get("x-auth-token")
		switch resp.StatusCode {
		case http.StatusOK:
			n = 11
		case http.StatusServiceUnavailable:
			//fmt.Printf("status code: %d request retryed\n", resp.StatusCode)
			time.Sleep(1 * time.Second)
			continue
		default:
			err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			return response, err
		}
	}
	return response, nil
}
