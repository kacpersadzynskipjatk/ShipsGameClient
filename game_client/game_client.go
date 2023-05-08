package game_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ServerAddress = "https://go-pjatk-server.fly.dev"
	GameEndpoint  = "/api/game"
	BoardEndpoint = "/api/game/board"
	FireEndpoint  = "/api/game/fire"
	DescEndpoint  = "/api/game/desc"
)

type GameClient struct {
	HttpClient http.Client
}

func NewGameClient(c *http.Client) *GameClient {
	gc := &GameClient{
		HttpClient: *c,
	}
	return gc
}

func (gc *GameClient) postRequest(address, token string, r *Request) (*http.Response, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, address, bytes.NewBuffer(jsonData))
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

func (gc *GameClient) PostStartGame(r Request) (StartGameResponse, error) {
	resp, err := gc.postRequest(ServerAddress+GameEndpoint, "", &r)
	if err != nil {
		return StartGameResponse{}, err
	}
	defer resp.Body.Close()

	startGameResponse := StartGameResponse{}
	err = startGameResponse.SetResponse(resp)
	if err != nil {
		return StartGameResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d, message: %s", resp.StatusCode, startGameResponse.Message)
		return StartGameResponse{}, err
	}
	return startGameResponse, nil
}

func (gc *GameClient) PostFire(token string, r Request) (FireResponse, error) {
	resp, err := gc.postRequest(ServerAddress+FireEndpoint, token, &r)
	if err != nil {
		return FireResponse{}, err
	}
	defer resp.Body.Close()

	fireResponse := FireResponse{}
	err = fireResponse.SetResponse(resp)
	if err != nil {
		return FireResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d, message: %s", resp.StatusCode, fireResponse.Message)
		return FireResponse{}, err
	}
	return fireResponse, nil
}

func (gc *GameClient) getRequest(address, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, address, nil)
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

func (gc *GameClient) GetGameStatus(token string) (StatusResponse, error) {
	resp, err := gc.getRequest(ServerAddress+GameEndpoint, token)
	if err != nil {
		return StatusResponse{}, err
	}
	defer resp.Body.Close()

	statusResponse := StatusResponse{}
	err = statusResponse.SetResponse(resp)
	if err != nil {
		return StatusResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d, message: %s", resp.StatusCode, statusResponse.Message)
		return StatusResponse{}, err
	}
	return statusResponse, nil
}

func (gc *GameClient) GetGameBoards(token string) (BoardResponse, error) {
	resp, err := gc.getRequest(ServerAddress+BoardEndpoint, token)
	if err != nil {
		return BoardResponse{}, err
	}
	defer resp.Body.Close()

	boardResponse := BoardResponse{}
	err = boardResponse.SetResponse(resp)
	if err != nil {
		return BoardResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d, message: %s", resp.StatusCode, boardResponse.Message)
		return BoardResponse{}, err
	}
	return boardResponse, nil
}

func (gc *GameClient) GetGameDescription(token string) (StatusResponse, error) {
	resp, err := gc.getRequest(ServerAddress+DescEndpoint, token)
	if err != nil {
		return StatusResponse{}, err
	}
	defer resp.Body.Close()

	descResponse := StatusResponse{}
	err = descResponse.SetResponse(resp)
	if err != nil {
		return StatusResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d, message: %s", resp.StatusCode, descResponse.Message)
		return StatusResponse{}, err
	}
	return descResponse, nil
}
