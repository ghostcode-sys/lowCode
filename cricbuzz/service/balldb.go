package service

import (
	"fmt"
	"lowleveldesign/cricbuzz/database"
	"time"
)

type ballDB struct {
	Id             int    `db:"id"`
	Ball_number    int    `db:"ball_number"`
	Bowler_id      string `db:"bowler_id"`
	Striker_id     string `db:"striker_id"`
	Non_striker_id string `db:"non_striker_id"`
	Runs_scored    int    `db:"runs_scored"`
	Is_wicket      bool   `db:"is_wicket"`
	Comment        string `db:"comment"`
	Extra_type     int `db:"extra_type"`
}

func init() {
	// To create a table for ball if not exists
	query := `CREATE TABLE IF NOT EXISTS ball (
			id SERIAL PRIMARY KEY,	
			ball_number INT NOT NULL,
			bowler_id VARCHAR(255) NOT NULL,
			striker_id VARCHAR(255) NOT NULL,
			non_striker_id VARCHAR(255) NOT NULL,
			runs_scored INT NOT NULL,
			is_wicket BOOLEAN NOT NULL,
			comment TEXT,
			extra_type INT
	)`

	conn := database.GetDBConnection()
	timeout := time.Second * 5

	err := database.CreateTable(conn, query, timeout)

	if err != nil {
		fmt.Println("Error creating a ball table: ", err.Error())
	}
}


func (b *ball) InsertBall() (int64, error) {
	conn := database.GetDBConnection()

	timeout := time.Second * 1

	query := `INSERT INTO ball (ball_number, bowler_id, striker_id, non_striker_id, runs_scored, is_wicket, comment, extra_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	params := []any{
		b.GetBallNumber(),
		b.GetBowler(),
		b.GetStriker(),
		b.GetNonStriker(),
		b.GetRunsScored(),
		b.GetWicket() != nil,
		b.GetComment(),
		b.GetExtra(),
	}

	rowAffected, err := database.ExecuteQuery(conn, query, params, timeout)

	return rowAffected, err
}
