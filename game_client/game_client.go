package game_client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	ServerAddress      = "https://go-pjatk-server.fly.dev/api/game"
	ServerBoardAddress = "https://go-pjatk-server.fly.dev/api/game/board"
	ServerFireAddress  = "https://go-pjatk-server.fly.dev/api/game/fire"
	ServerDescAddress  = "https://go-pjatk-server.fly.dev/api/game/desc"
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

func (gc *GameClient) postRequest(address, token string, r *Request) *http.Response {
	jsonData, _ := json.Marshal(r)
	req, _ := http.NewRequest("POST", address, bytes.NewBuffer(jsonData))
	req.Header.Add("X-auth-token", token)
	resp, err := gc.HttpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (gc *GameClient) PostStartGame(r Request) StartGameResponse {
	resp := gc.postRequest(ServerAddress, "", &r)
	defer resp.Body.Close()

	startGameResponse := StartGameResponse{}
	startGameResponse.SetResponse(resp)
	return startGameResponse
}

func (gc *GameClient) PostFire(token string, r Request) FireResponse {
	resp := gc.postRequest(ServerFireAddress, token, &r)
	defer resp.Body.Close()

	fireResponse := FireResponse{}
	fireResponse.SetResponse(resp)
	return fireResponse
}

func (gc *GameClient) getRequest(address, token string) *http.Response {
	req, _ := http.NewRequest("GET", address, nil)
	req.Header.Add("X-auth-token", token)
	resp, err := gc.HttpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (gc *GameClient) GetGameStatus(token string) StatusResponse {
	resp := gc.getRequest(ServerAddress, token)
	defer resp.Body.Close()

	statusResponse := StatusResponse{}
	statusResponse.SetResponse(resp)
	return statusResponse
}

func (gc *GameClient) GetGameBoards(token string) BoardResponse {
	resp := gc.getRequest(ServerBoardAddress, token)
	defer resp.Body.Close()

	boardResponse := BoardResponse{}
	boardResponse.SetResponse(resp)
	return boardResponse
}

func (gc *GameClient) GetGameDescription(token string) StatusResponse {
	resp := gc.getRequest(ServerDescAddress, token)
	defer resp.Body.Close()

	descResponse := StatusResponse{}
	descResponse.SetResponse(resp)
	return descResponse
}
