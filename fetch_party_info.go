package main

import (
	"encoding/json"
	"net/http"
)

const partyInfoAPI = "https://raw.githubusercontent.com/basputtipong/election-data/main/data/info_party_overview.json"

type PartyInfoResponse struct {
	ID      string `json:"id"`
	PartyNo string `json:"party_no"`
	Name    string `json:"name"`
}

func FetchPartyInfo() ([]PartyInfoResponse, error) {
	resp, err := http.Get(partyInfoAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []PartyInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, err
}
