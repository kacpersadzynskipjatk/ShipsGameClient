package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

const ServerAddress = "https://go-pjatk-server.fly.dev/api"

// Game endpoints
const (
	GameEndpoint    = "/game"
	BoardEndpoint   = "/game/board"
	FireEndpoint    = "/game/fire"
	DescEndpoint    = "/game/desc"
	AbandonEndpoint = "/game/abandon"
	RefreshEndpoint = "/game/refresh"
)

// Lobby endpoint
const LobbyEndpoint = "/lobby"

// Player stats endpoint
const PlayerStatsEndpoint = "/stats"

// Client represents a client for making HTTP requests with an associated token.
type Client struct {
	HttpClient http.Client // The underlying HTTP client for making requests
	Token      string      // The authentication token to be used in requests
}

// NewClient creates a new instance of the Client struct.
// It takes a pointer to a http.Client as a parameter and returns a pointer to the newly created Client.
func NewClient(client *http.Client) *Client {
	newClient := &Client{
		HttpClient: *client, // Assign the provided http.Client to the HttpClient field
	}
	return newClient
}

// doRequest performs an HTTP request based on the provided request type, address, token, and payload.
// It returns the HTTP response and any error that occurred during the request.
func (c *Client) doRequest(requestType, address, token string, payload interface{}) (*http.Response, error) {
	var req *http.Request
	var err error

	// Create the HTTP request based on the provided parameters
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

	// Add the token to the request header
	req.Header.Add("X-auth-token", token)

	// Perform the HTTP request using the client's underlying HTTP client
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SendRequest sends an HTTP request to the specified endpoint using the provided method type, token, request payload, and response struct.
// It retries the request up to 10 times in case of specific errors, and returns the response or any error that occurred.
func (c *Client) SendRequest(methodType, endpoint, token string, request interface{}, response interface{}) (interface{}, error) {
	for n := 1; n <= 10; n++ {
		var resp *http.Response
		var err error

		// Determine the request payload and call the doRequest function to perform the HTTP request
		if request != nil {
			resp, err = c.doRequest(methodType, ServerAddress+endpoint, token, request)
		} else {
			resp, err = c.doRequest(methodType, ServerAddress+endpoint, token, nil)
		}

		if err != nil {
			return response, err
		}
		defer resp.Body.Close()

		// Decode the response into the provided response struct
		if response != nil {
			err = json.NewDecoder(resp.Body).Decode(response)
		}

		if err != nil {
			if err != io.EOF {
				return response, err
			}
		}

		// Update the client's token based on the response header
		c.Token = resp.Header.Get("x-auth-token")

		// Handle specific HTTP status codes and retry or return errors accordingly
		switch resp.StatusCode {
		case http.StatusOK:
			n = 11
		case http.StatusServiceUnavailable:
			time.Sleep(1 * time.Second)
			continue
		case http.StatusNotFound:
			err = fmt.Errorf("nie znaleziono, kod błędu:%d\n", resp.StatusCode)
			return response, err
		default:
			err = fmt.Errorf("unexpected status code:%d error: %w\n%s", resp.StatusCode, err, debug.Stack())
			return response, err
		}
	}
	return response, nil
}
