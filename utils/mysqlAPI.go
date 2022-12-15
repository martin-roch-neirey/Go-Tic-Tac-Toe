package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

func getDatabaseConnection() *sql.DB {
	db, err := sql.Open("mysql", "gouser:pass@tcp(127.0.0.1:3306)/tictactoe")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetGamesCount() int {
	db := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	row := db.QueryRow("SELECT COUNT(*) FROM games")

	var count int
	err := row.Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	return count
}

func UploadNewGame(gamemode string, winner int) {
	db := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "INSERT INTO games(gamemode, winner) VALUES('VAL1',VAL2);"
	query = strings.Replace(query, "VAL1", gamemode, 1)
	query = strings.Replace(query, "VAL2", strconv.Itoa(winner), 1)

	fmt.Println(query)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func closeDatabaseConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		return
	}
}

func main() {
	UploadNewGame("test", 91)
	fmt.Println(GetGamesCount())

	UploadNewGame("test2", 863)
	fmt.Println(GetGamesCount())
}
