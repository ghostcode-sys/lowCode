package composer

import (
	"encoding/json"
	"lowleveldesign/cricbuzz/service"
	"os"
	"strings"
)

func CreatePlayers() ([]*service.PlayerInterface, error) {
	playerFile := "./records/player.json"
	players, err := os.ReadFile(playerFile)

	playerList := make([]*service.PlayerInterface, 0)

	if err != nil {
		return playerList, err
	}

	type playerStruct struct {
		PlayerId     string `json:"playerId"`
		PlayerName   string `json:"playerName"`
		PlayerRole   string `json:"playerRole"`
		BattingStyle string `json:"battingStyle"`
		BowlingStyle string `json:"bowlingStyle"`
		Country      string `json:"country"`
	}

	playerArr := make([]playerStruct, 0)
	err = json.Unmarshal(players, &playerArr)
	if err != nil {
		return playerList, err
	}

	for _, p := range playerArr {
		var player service.PlayerInterface
		player = service.NewPlayer(p.PlayerName, p.PlayerId)

		switch p.PlayerRole {
		case "Batsman":
			player.SetPlayerRole(service.BATSMAN)
		case "Bowler":
			player.SetPlayerRole(service.BOWLER)
		case "Wicket Keeper":
			player.SetPlayerRole(service.WICKET_KEEPER)
		case "All Rounder":
			player.SetPlayerRole(service.ALL_ROUNDER)
		}

		player.SetCountry(p.Country)

		switch p.BattingStyle {
		case "Right Handed":
			player.SetBattingStyle(service.RIGHTHANDY)
		case "Left Handed":
			player.SetBattingStyle(service.LEFTHANDY)
		}

		bowlDetails := strings.Split(p.BowlingStyle, " ")

		dominateHand := service.RIGHTHANDY

		if len(bowlDetails) > 0 {
			switch bowlDetails[0] {
			case "Right-arm":
				dominateHand = service.RIGHTHANDY
			case "Left-arm":
				dominateHand = service.LEFTHANDY
			}
		}

		bowlType := ""
		if len(bowlDetails) > 1 {
			bowlType = strings.Join(bowlDetails[1:], " ")
		}

		player.SetBowlingStyle(bowlType, dominateHand)

		playerList = append(playerList, &player)
	}
	return playerList, nil
}
