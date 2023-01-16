package api

import (
	"database/sql"
	"errors"
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

func getDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")
	if err != nil {
		return nil, errors.New("mySQL DB not reachable")
	}
	return db, nil
}

func IsSqlApiUsable() bool {
	db, err := getDatabaseConnection()
	pingError := db.Ping()
	if pingError != nil || err != nil {
		return false
	}
	return true
}

func GetGamesCount() int {
	if isGamesCountCacheValid {
		return gamesCountCache
	}

	db, _ := getDatabaseConnection()
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

	db, _ := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "SELECT * FROM (SELECT * FROM games3 ORDER BY id DESC LIMIT 5) AS sub ORDER BY id DESC;"
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
	db, _ := getDatabaseConnection()
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
