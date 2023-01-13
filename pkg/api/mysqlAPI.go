package api

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

var isGamesCountCacheValid = false
var gamesCountCache int

var isLastGamesCacheValid = false
var lastGamesCache []string

func getDatabaseConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetGamesCount() int {
	if isGamesCountCacheValid {
		return gamesCountCache
	}

	db := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	row := db.QueryRow("SELECT COUNT(*) FROM games3")

	var count int
	err := row.Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	gamesCountCache = count
	go validateCache(&isGamesCountCacheValid)
	return count
}

func GetLastGames() []string {

	if isLastGamesCacheValid {
		return lastGamesCache
	}

	db := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "SELECT * FROM (SELECT * FROM games3 ORDER BY id DESC LIMIT 5) AS sub ORDER BY id ASC;"
	// query = strings.Replace(query, "VAL1", strconv.Itoa(number), 1)

	var games []string
	rows, _ := db.Query(query)

	for rows.Next() {
		var value string
		var id int
		var date string
		err := rows.Scan(&id, &date, &value)
		if err != nil {
			log.Fatal(err)
		} else {
			games = append(games, value)
		}
	}

	fmt.Println(games)

	go validateCache(&isLastGamesCacheValid)
	lastGamesCache = games
	return games
}

func UploadNewGame(json string) {
	db := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "INSERT INTO games3(date, properties) VALUES('VAL1', 'VAL2');"
	query = strings.Replace(query, "VAL1", time.Now().Format("2006-01-02 15-04-05"), 1)
	query = strings.Replace(query, "VAL2", json, 1)

	// fmt.Println(query)

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

func validateCache(variable *bool) {
	*variable = true
	time.Sleep(10 * time.Second)
	*variable = false
}
func main() { // to test, change package to main in this file and all files of the folder utils
	UploadNewGame("{}")
	/*
		fmt.Println(GetGamesCount())

		UploadNewGame("{}")
		fmt.Println(GetGamesCount())*/

	GetLastGames()
}
