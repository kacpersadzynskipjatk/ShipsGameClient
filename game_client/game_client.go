package game_client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	ServerAddress = "https://go-pjatk-server.fly.dev/api/game"
	ServerBoardAddress = "https://go-pjatk-server.fly.dev/api/game/board"
	ServerFireAddress = "https://go-pjatk-server.fly.dev/api/game/fire"
)

type GameClient struct {
	HttpClient http.Client
}

func NewGameClient(c *http.Client) (*GameClient) {
	gc := &GameClient{
		HttpClient: *c,
	}
	return gc
}

func (gc *GameClient) PostStartGame(params map[string]any) string{
	var reqBody []byte = nil
	if params != nil {
		reqBody, _ = json.Marshal(params)
	}
	req, _ := http.NewRequest("POST", ServerAddress, bytes.NewBuffer(reqBody))
	resp, _ := gc.HttpClient.Do(req)
	defer resp.Body.Close()
	return resp.Header.Get("x-auth-token")
}

func (gc *GameClient) GetGameStatus(){

}

func (gc *GameClient) GetGameBoards(){

}

func (gc *GameClient) PostFire(){

}
