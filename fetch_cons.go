package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var partyConsAPI = "https://raw.githubusercontent.com/basputtipong/election-data/main/data/stats_cons.json"

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
	req, _ := http.NewRequest("GET", partyConsAPI, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("fetch failed: status=%d body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("non-200 response")
	}

	var data ElectionResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("json unmarshal error, body=%s", string(body))
		return nil, err
	}
	return &data, err
}
