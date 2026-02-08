package main

import (
	"encoding/json"
	"net/http"
)

type ProvinceInfos struct {
	TotalRegisterVote int            `json:"total_registered_vote"`
	TotalStationVote  int            `json:"total_vote_stations"`
	Province          []ProvinceInfo `json:"province"`
}

type ProvinceInfo struct {
	ProvinceID string `json:"province_id"`
	CityCode   string `json:"prov_id"`
	Name       string `json:"province"`
}

var provinceAPI = "https://raw.githubusercontent.com/basputtipong/election-data/main/data/info_province.json"

func FetchProvinceInfo() (*ProvinceInfos, error) {
	resp, err := http.Get(provinceAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ProvinceInfos
	err = json.NewDecoder(resp.Body).Decode(&data)
	return &data, err
}
