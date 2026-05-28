package service

import (
	"fmt"
	"lowleveldesign/cricbuzz/database"
	"time"
)

type playerDB struct {
	ID                  int          `db:"id"`
	PlayerID            string       `db:"player_id"`
	PlayerName          string       `db:"player_name"`
	PlayerRole          PlayerRole   `db:"player_role"`
	BattingStyle        DominateSide `db:"batting_style"`
	BowlingDominateSide DominateSide `db:"bowling_dominate_side"`
	BowlingType         string       `db:"bowling_type"`
	Country             string       `db:"country"`
}

func init() {
	// To create a table for player if not exists
	query := `CREATE TABLE IF NOT EXISTS player (
			id SERIAL PRIMARY KEY,	
			player_id VARCHAR(255) UNIQUE NOT NULL,
			player_name VARCHAR(255) NOT NULL,
			player_role VARCHAR(50) NOT NULL,
			batting_style VARCHAR(50) NOT NULL,
			bowling_style VARCHAR(50),
			country VARCHAR(100) NOT NULL
	)`

	conn := database.GetDBConnection()

	timeout := time.Second * 5

	err := database.CreateTable(conn, query, timeout)

	if err != nil {
		fmt.Println("Error creating a player table: ", err.Error())
	}
}

func (p *player) InsertPlayer() (int64, error) {
	conn := database.GetDBConnection()

	timeout := time.Second * 1

	query := `INSERT INTO player (player_id, player_name, player_role, batting_style, bowling_style, country) VALUES (?, ?, ?, ?, ?, ?, ?)`

	params := []any{
		p.GetPlayerID(),
		p.GetPlayerName(),
		p.GetPlayerRole(),
		p.GetBattingStyle(),
		p.GetBowlingStyle(),
		p.GetCountry(),
	}

	rowAffected, err := database.ExecuteQuery(conn, query, params, timeout)

	if err != nil {
		fmt.Printf("Error while creating Player: %s \n", err.Error())
	}

	return rowAffected, err
}

func (p *player) UpdatePlayer() (int64, error) {
	conn := database.GetDBConnection()

	timeout := time.Second * 1

	query := `UPDATE player SET player_name = ?, player_role = ?, batting_style = ?, bowling_style = ?, country = ? WHERE player_id = ?`

	params := []any{
		p.GetPlayerName(),
		p.GetPlayerRole(),
		p.GetBattingStyle(),
		p.GetBowlingStyle(),
		p.GetCountry(),
		p.GetPlayerID(),
	}

	rowAffected, err := database.ExecuteQuery(conn, query, params, timeout)

	if err != nil {
		fmt.Printf("Error while updating Player: %s \n", err.Error())
	}

	return rowAffected, err
}

func (p *player) GetPlayerInfo() (*playerDB, error) {
	conn := database.GetDBConnection()

	timeout := time.Second * 1

	query := `SELECT id, player_id, player_name, player_role, batting_style, bowling_style, country FROM player WHERE player_id = ?`

	params := []any{
		p.GetPlayerID(),
	}

	var playerInfo playerDB

	res, err := database.SelectQuery[playerDB](conn, query, params, timeout)

	if err != nil {
		fmt.Printf("Error while fetching Player info: %s \n", err.Error())
		return nil, err
	}

	if len(res) > 0 {
		playerInfo = res[0]
	}

	return &playerInfo, nil
}
