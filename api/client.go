package api

import (
	"encoding/json"
	"go-cs-log-parser/database"
	"io"
	"net/http"
	"strconv"
	"time"
)

func GetPlayers(apiURL string) ([]database.ListPlayersRow, error) {
	apiClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, apiURL+"players", nil)
	if err != nil {
		return nil, err
	}
	res, getErr := apiClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	var players []database.ListPlayersRow
	jsonErr := json.Unmarshal(body, &players)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return players, nil
}

func GetPlayerWithStats(apiURL string, playerID int64) (*database.GetPlayerWithStatsRow, error) {
	apiClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, apiURL+"players/"+strconv.FormatInt(playerID, 10), nil)
	if err != nil {
		return nil, err
	}
	res, getErr := apiClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	var playerWithStats database.GetPlayerWithStatsRow
	jsonErr := json.Unmarshal(body, &playerWithStats)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &playerWithStats, nil
}

func GetAccoladeForPlayer(apiURL string, playerID int64) ([]database.GetAccoladeForPlayerRow, error) {
	apiClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, apiURL+"players/"+strconv.FormatInt(playerID, 10)+"/accolade", nil)
	if err != nil {
		return nil, err
	}
	res, getErr := apiClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	var playerAccolade []database.GetAccoladeForPlayerRow
	jsonErr := json.Unmarshal(body, &playerAccolade)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return playerAccolade, nil
}
