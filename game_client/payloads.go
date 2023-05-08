package game_client

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response interface {
	SetResponse(*http.Response) error
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

	Message string `json:"message"`
}

func (r *StatusResponse) SetResponse(rawData *http.Response) error {
	err := json.NewDecoder(rawData.Body).Decode(&r)
	if err != nil {
		return err
	}
	return nil
}

type BoardResponse struct {
	Board []string `json:"board"`

	Message string `json:"message"`
}

func (r *BoardResponse) SetResponse(rawData *http.Response) error {
	err := json.NewDecoder(rawData.Body).Decode(&r)
	if err != nil {
		return err
	}
	return nil
}

type StartGameResponse struct {
	Token string

	Message string `json:"message"`
}

func (r *StartGameResponse) SetResponse(rawData *http.Response) error {
	r.Token = rawData.Header.Get("x-auth-token")
	if r.Token == "" {
		return errors.New("token is empty")
	}
	return nil
}

type FireResponse struct {
	Result string `json:"result"`

	Message string `json:"message"`
}

func (r *FireResponse) SetResponse(rawData *http.Response) error {
	err := json.NewDecoder(rawData.Body).Decode(&r)
	if err != nil {
		return err
	}
	return nil
}

type Request interface {
}

type FireRequest struct {
	Coord string `json:"coord"`
}

type StartGameRequest struct {
	Coords     []string `json:"coords"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick"`
	Wpbot      bool     `json:"wpbot"`
}
