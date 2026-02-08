package main

import (
	"encoding/json"
	"net/http"
)

const partyAPI = "https://stats-ectreport69.ect.go.th/data/records/stats_party.json"

type PartyResponse struct {
	CountedVoteStations int           `json:"counted_vote_stations"`
	PercentCount        float64       `json:"percent_count"`
	ResultParty         []PartyResult `json:"result_party"`
}

type PartyResult struct {
	PartyID          int         `json:"party_id"`
	PartyVote        int         `json:"party_vote"`
	PartyVotePercent float64     `json:"party_vote_percent"`
	PartyListCount   *int        `json:"party_list_count"`
	MPAppVote        int         `json:"mp_app_vote"`
	MPAppVotePercent float64     `json:"mp_app_vote_percent"`
	FirstMPAppCount  int         `json:"first_mp_app_count"`
	Candidates       []Candidate `json:"candidates"`
}

type Candidate struct {
	MPAppID          string  `json:"mp_app_id"`
	MPAppVote        int     `json:"mp_app_vote"`
	MPAppVotePercent float64 `json:"mp_app_vote_percent"`
	MPAppRank        int     `json:"mp_app_rank"`
	PartyID          int     `json:"party_id"`
}

func FetchParty() (*PartyResponse, error) {
	resp, err := http.Get(partyAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data PartyResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return &data, err
}
