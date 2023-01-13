package api

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
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

func GetLastGames(number int) {
	db := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "SELECT * FROM (SELECT * FROM games3 ORDER BY id DESC LIMIT 3) AS sub ORDER BY id ASC;"
	query = strings.Replace(query, "VAL1", strconv.Itoa(number), 1)

	var games []string
	rows, _ := db.Query(query)

	for rows.Next() {
		var value string
		err := rows.Scan(value)
		if err != nil {
			log.Fatal(err)
		} else {
			games = append(games, value)
		}
	}

	fmt.Println(games)
}

func UploadNewGame(json string) {
	db := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "INSERT INTO games3(date, properties) VALUES('VAL1', 'VAL2');"
	query = strings.Replace(query, "VAL1", time.Now().Format("2006-01-02 15-04-05"), 1)
	query = strings.Replace(query, "VAL2", json, 1)

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

func main() { // to test, change package to main in this file and all files of the folder utils
	UploadNewGame("{}")
	/*
		fmt.Println(GetGamesCount())

		UploadNewGame("{}")
		fmt.Println(GetGamesCount())*/

	GetLastGames(3)
}
