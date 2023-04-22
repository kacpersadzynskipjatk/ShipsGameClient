package game_client

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	ServerAddress      = "https://go-pjatk-server.fly.dev/api/game"
	ServerBoardAddress = "https://go-pjatk-server.fly.dev/api/game/board"
	ServerFireAddress  = "https://go-pjatk-server.fly.dev/api/game/fire"
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

func (gc *GameClient) PostStartGame(params map[string]any) string {

	var reqBody []byte = nil
	if params != nil {
		reqBody, _ = json.Marshal(params)
	}

	req, _ := http.NewRequest("POST", ServerAddress, bytes.NewBuffer(reqBody))
	resp, err := gc.HttpClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	return resp.Header.Get("x-auth-token")
}

func (gc *GameClient) PostFire(coord, token string) string {

	data := map[string]string{"coord": coord}
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", ServerFireAddress, bytes.NewBuffer(jsonData))
	req.Header.Add("X-auth-token", token)
	resp, err := gc.HttpClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	response := string(body)
	return response
}

func (gc *GameClient) GetRequest(address, token string) string {
	req, _ := http.NewRequest("GET", address, nil)
	req.Header.Add("X-auth-token", token)
	resp, err := gc.HttpClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	response := string(body)
	return response
}

func (gc *GameClient) GetGameStatus(token string) string {
	return gc.GetRequest(ServerAddress, token)
}

func (gc *GameClient) GetGameBoards(token string) string {
	return gc.GetRequest(ServerBoardAddress, token)
}
