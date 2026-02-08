package main

import (
	"encoding/json"
	"net/http"
)

var partyConsAPI = "https://stats-ectreport69.ect.go.th/data/records/stats_cons.json"

type ElectionResponse struct {
	ResultProvince []Province `json:"result_province"`
	LastUpdate     string     `json:"last_update"`
}

type Province struct {
	ProvinceID      string            `json:"prov_id"`
	PartyResultCons []PartyResultCons `json:"result_party"`
}

type PartyResultCons struct {
	PartyID              int     `json:"party_id"`
	PartyListVote        int     `json:"party_list_vote"`
	PartyListVotePercent float64 `json:"party_list_vote_percent"`
}

func FetchCons() (*ElectionResponse, error) {
	resp, err := http.Get(partyConsAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ElectionResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return &data, err
}
