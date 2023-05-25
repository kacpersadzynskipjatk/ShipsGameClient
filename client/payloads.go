package client

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

type BoardResponse struct {
	Board []string `json:"board"`

	Message string `json:"message"`
}

type StartGameResponse struct {
	Token string

	Message string `json:"message"`
}

type FireResponse struct {
	Result string `json:"result"`

	Message string `json:"message"`
}

type AbandonResponse struct {
	Message string `json:"message"`
}

type OpponentResponse struct {
	GameStatus string `json:"game_status"`
	Nick       string `json:"nick"`

	Message string `json:"message"`
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
