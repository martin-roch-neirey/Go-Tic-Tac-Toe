// Copyright (c) 2022 Haute école d'ingénierie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

package api

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

var isGamesCountCacheValid = false
var gamesCountCache int

var isLastGamesCacheValid = false
var lastGamesCache []string

// getDatabaseConnection returns a sql.DB object
// that represents connection to mySQL database
func getDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:test@tcp(51.210.255.160:3306)/tictactoe")
	if err != nil {
		return nil, errors.New("mySQL DB not reachable")
	}
	return db, nil
}

// IsSqlApiUsable returns true if mySQL database is reachable, false otherwise
func IsSqlApiUsable() bool {
	db, err := getDatabaseConnection()
	pingError := db.Ping()
	if pingError != nil || err != nil {
		return false
	}
	return true
}

// GetGamesCount returns number of rows in games table
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
	go validateCache(&isGamesCountCacheValid) // validate cache after sql query
	return count
}

// GetLastGames returns the five last rows inserted to games table
func GetLastGames() []string {

	if isLastGamesCacheValid {
		return lastGamesCache
	}

	db, _ := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "SELECT * FROM (SELECT * FROM games3 ORDER BY id DESC LIMIT 5) AS sub ORDER BY id DESC;"

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

	go validateCache(&isLastGamesCacheValid) // validate cache after sql query
	lastGamesCache = games
	return games
}

// UploadNewGame uploads given json to database as a new game entry
func UploadNewGame(json string) {
	db, _ := getDatabaseConnection()
	defer closeDatabaseConnection(db)

	var query string
	query = "INSERT INTO games3(date, properties) VALUES('VAL1', 'VAL2');"
	query = strings.Replace(query, "VAL1", time.Now().Format("2006-01-02 15-04-05"), 1)
	query = strings.Replace(query, "VAL2", json, 1)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// closeDatabaseConnection closes db connection
func closeDatabaseConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		return
	}
}

// validateCache is a function used to validate cache of a specific object type.
// Cache will be valid 10 seconds
func validateCache(variable *bool) {
	*variable = true
	time.Sleep(10 * time.Second)
	*variable = false
}
