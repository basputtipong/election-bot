package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func BuildElectionEmbed(
	party *PartyResponse,
	partyInfo []PartyInfoResponse,
	partyCons *ElectionResponse,
	provinceInfo *ProvinceInfos,
) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: "🗳️ ผลการนับคะแนนเลือกตั้ง (Realtime)",
		Color: 0xE91E63,
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "📊 ความคืบหน้า",
		Value: fmt.Sprintf("นับแล้ว %.2f%% (%s หน่วย)", party.PercentCount, formatNumber(party.CountedVoteStations)),
	})

	if len(party.ResultParty) == 0 {
		embed.Description = "⚠️ ยังไม่มีข้อมูลพรรค"
		return embed
	}

	sort.Slice(party.ResultParty, func(i, j int) bool {
		return party.ResultParty[i].PartyVote > party.ResultParty[j].PartyVote
	})

	limit := 5
	if len(party.ResultParty) < limit {
		limit = len(party.ResultParty)
	}

	mapPartyName := setPartyName(partyInfo)

	for i := 0; i < limit; i++ {
		p := party.ResultParty[i]
		name := mapPartyName[strconv.Itoa(p.PartyID)]
		emoji := rankEmoji(i + 1)

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s 🏛️ พรรค %s", emoji, name),
			Value:  fmt.Sprintf("คะแนน: %s (%.2f%%)", formatNumber(p.PartyVote), p.PartyVotePercent),
			Inline: false,
		})
	}

	return embed
}

func BuildElectionCityEmbed(
	party *PartyResponse,
	partyInfo []PartyInfoResponse,
	partyCons *ElectionResponse,
	provinceInfo *ProvinceInfos,
	provinceName string,
) *discordgo.MessageEmbed {
	provinceCode := getProvinceIDByName(provinceName, provinceInfo)

	lastUpdate, _ := time.Parse(time.RFC3339Nano, partyCons.LastUpdate)
	timeZone, _ := time.LoadLocation("Asia/Bangkok")
	lastUpdate = lastUpdate.In(timeZone)

	cityEmbed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("🗳️ ผลการนับคะแนนเลือกตั้งจังหวัด: %s", provinceName),
		Color: 0xE91E63,
	}

	if provinceCode == "" {
		cityEmbed.Description = fmt.Sprintf("⚠️ ไม่มีข้อมูลของจังหวัด: %s", provinceName)
		return cityEmbed
	}

	cityEmbed.Fields = append(cityEmbed.Fields, &discordgo.MessageEmbedField{
		Name:  "⏱️ อัพเดทล่าสุด",
		Value: fmt.Sprintf("เวลา: %s", lastUpdate),
	})

	provinceList := make(map[string][]PartyResultCons, 0)
	for _, prov := range partyCons.ResultProvince {
		sort.Slice(prov.PartyResultCons, func(i, j int) bool {
			return prov.PartyResultCons[i].PartyListVote > prov.PartyResultCons[j].PartyListVote
		})

		provinceList[prov.ProvinceID] = prov.PartyResultCons
	}

	limit := 5
	for i := 0; i < limit; i++ {
		mapPartyName := setPartyName(partyInfo)
		partyName := mapPartyName[strconv.Itoa(provinceList[provinceCode][i].PartyID)]
		emoji := rankEmoji(i + 1)
		cityEmbed.Fields = append(cityEmbed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s 🏛️ พรรค %s", emoji, partyName),
			Value:  fmt.Sprintf("คะแนน: %s (%.2f%%)", formatNumber(provinceList[provinceCode][i].PartyListVote), provinceList[provinceCode][i].PartyListVotePercent),
			Inline: false,
		})
	}

	return cityEmbed
}

func setPartyName(partyInfo []PartyInfoResponse) map[string]string {
	mapPartyNames := make(map[string]string)
	for _, p := range partyInfo {
		mapPartyNames[p.ID] = p.Name
	}

	return mapPartyNames
}

func formatNumber(n int) string {
	s := strconv.Itoa(n)
	nLen := len(s)

	if nLen <= 3 {
		return s
	}

	var result []byte
	pre := nLen % 3
	if pre > 0 {
		result = append(result, s[:pre]...)
		if nLen > pre {
			result = append(result, ',')
		}
	}

	for i := pre; i < nLen; i += 3 {
		result = append(result, s[i:i+3]...)
		if i+3 < nLen {
			result = append(result, ',')
		}
	}

	return string(result)
}

func rankEmoji(rank int) string {
	switch rank {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		return ""
	}
}

func getProvinceIDByName(provinceName string, provinceInfo *ProvinceInfos) string {
	provinceInfoMap := make(map[string]string)
	for _, info := range provinceInfo.Province {
		provinceInfoMap[info.Name] = info.CityCode
	}
	return provinceInfoMap[provinceName]
}
