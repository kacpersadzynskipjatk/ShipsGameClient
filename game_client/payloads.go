package game_client

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

type Response interface {
	SetResponse(*http.Response)
	GetResponse() Response
}

type StatusResponse struct {
	Desc           string   `json:"desc"`
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Nick           string   `json:"nick"`
	OppDesc        string   `json:"opp_desc"`
	OppShots       []string `json:"opp_shots"`
	Opponent       string   `json:"opponent"`
	ShouldFire     bool     `json:"should_fire"`
	Timer          int      `json:"timer"`
}

func (r *StatusResponse) SetResponse(rawData *http.Response) {
	json.NewDecoder(rawData.Body).Decode(&r)
}

func (r *StatusResponse) GetResponse() Response {
	return r
}

type BoardResponse struct {
	Board []string `json:"board"`
}

func (r *BoardResponse) SetResponse(rawData *http.Response) {
	json.NewDecoder(rawData.Body).Decode(&r)
}

func (r *BoardResponse) GetResponse() Response {
	return r
}

type StartGameResponse struct {
	Token   string
	Message string `json:"message"`
}

func (r *StartGameResponse) SetResponse(rawData *http.Response) {
	if rawData.StatusCode != http.StatusOK {
		json.NewDecoder(rawData.Body).Decode(&r)
	} else {
		r.Token = rawData.Header.Get("x-auth-token")
	}
}

func (r *StartGameResponse) GetResponse() Response {
	return r
}

type FireResponse struct {
	Result string `json:"result"`
}

func (r *FireResponse) SetResponse(rawData *http.Response) {
	json.NewDecoder(rawData.Body).Decode(&r)
}

func (r *FireResponse) GetResponse() Response {
	return r
}

type Request interface {
	SetRequest(map[string]any)
	GetRequest() Request
}

type FireRequest struct {
	Coord string `json:"coord"`
}

func (r *FireRequest) SetRequest(rawData map[string]any) {
	mapstructure.Decode(rawData, &r)
}

func (r *FireRequest) GetRequest() Request {
	return r
}

type StartGameRequest struct {
	Coords     []string `json:"coords"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick"`
	Wpbot      bool     `json:"wpbot"`
}

func (r *StartGameRequest) SetRequest(rawData map[string]any) {
	mapstructure.Decode(rawData, &r)
}

func (r *StartGameRequest) GetRequest() Request {
	return r
}
