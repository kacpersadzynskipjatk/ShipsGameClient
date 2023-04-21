package game_client

import "net/http"

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

func (gc *GameClient) PostStartGame(){
	req, err := http.NewRequest("POST", ServerAddress, nil)
}

func (gc *GameClient) GetGameStatus(){

}

func (gc *GameClient) GetGameBoards(){

}

func (gc *GameClient) PostFire(){

}
