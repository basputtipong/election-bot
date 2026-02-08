package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	lastMessageIDMain string
	lastMessageIDCity string
	targetChannel     string
	mu                sync.Mutex
)

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	mu.Lock()
	targetChannel = m.ChannelID

	if strings.EqualFold(m.Content, "!election") {
		lastMessageIDMain = ""
		mu.Unlock()

		log.Println("!election command received")

		sendOrUpdateElection(s)
	} else if strings.HasPrefix(strings.ToLower(m.Content), "!election ") {
		lastMessageIDCity = ""
		mu.Unlock()
		parts := strings.Fields(m.Content)

		var provinceName string
		if len(parts) > 1 {
			provinceName = strings.Join(parts[1:], " ")
		}

		sendElection(s, provinceName)

		fmt.Println("!election command received:", provinceName)
	}
}

func sendOrUpdateElection(s *discordgo.Session) {
	mu.Lock()
	channelID := targetChannel
	mainID := lastMessageIDMain
	mu.Unlock()

	if channelID == "" {
		return
	}

	partyCons, err := FetchCons()
	if err != nil {
		log.Println("fetch party cons failed:", err)
		return
	}

	provinceInfo, err := FetchProvinceInfo()
	if err != nil {
		log.Println("fetch province info failed:", err)
		return
	}

	party, err := FetchParty()
	if err != nil {
		log.Println("fetch party failed:", err)
		return
	}

	partyInfo, err := FetchPartyInfo()
	if err != nil {
		log.Println("fetch party info failed:", err)
		return
	}

	embedMain := BuildElectionEmbed(
		party,
		partyInfo,
		partyCons,
		provinceInfo,
	)

	if mainID != "" {
		_, err = s.ChannelMessageEditEmbed(channelID, mainID, embedMain)
		if err != nil {
			log.Println("edit main embed failed:", err)
		}
		return
	}

	msgMain, err := s.ChannelMessageSendEmbed(channelID, embedMain)
	if err != nil {
		log.Println("send main embed failed:", err)
		return
	}

	mu.Lock()
	lastMessageIDMain = msgMain.ID
	mu.Unlock()
}

func sendElection(s *discordgo.Session, provinceName string) {
	mu.Lock()
	channelID := targetChannel
	cityID := lastMessageIDCity
	mu.Unlock()

	if channelID == "" {
		return
	}

	partyCons, err := FetchCons()
	if err != nil {
		log.Println("fetch party cons failed:", err)
		return
	}

	provinceInfo, err := FetchProvinceInfo()
	if err != nil {
		log.Println("fetch province info failed:", err)
		return
	}

	party, err := FetchParty()
	if err != nil {
		log.Println("fetch party failed:", err)
		return
	}

	partyInfo, err := FetchPartyInfo()
	if err != nil {
		log.Println("fetch party info failed:", err)
		return
	}

	embedCity := BuildElectionCityEmbed(
		party,
		partyInfo,
		partyCons,
		provinceInfo,
		provinceName,
	)

	if cityID != "" {
		_, err = s.ChannelMessageEditEmbed(channelID, cityID, embedCity)
		if err != nil {
			log.Println("edit main embed failed:", err)
		}
		return
	}

	msgCity, err := s.ChannelMessageSendEmbed(channelID, embedCity)
	if err != nil {
		log.Println("send main embed failed:", err)
		return
	}

	mu.Lock()
	lastMessageIDCity = msgCity.ID
	mu.Unlock()
}
