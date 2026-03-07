package db

import (
	"database/sql"
	"log"
	"masjid_baiturrahman/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./masjid.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS prayers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		time TEXT
	);
	`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// Check if data exists
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM prayers").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		insert := `
		INSERT INTO prayers (name, time) VALUES
		('Subuh', '05:12'),
		('Syuruq', '06:26'),
		('Dzuhur', '12:34'),
		('Ashar', '15:52'),
		('Maghrib', '18:38'),
		('Isya', '19:48');
		`
		_, err = DB.Exec(insert)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetPrayers() ([]models.Prayer, error) {
	rows, err := DB.Query("SELECT id, name, time FROM prayers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prayers []models.Prayer
	now := time.Now()
	currentTime := now.Format("15:04")
	nextFound := false

	for rows.Next() {
		var p models.Prayer
		if err := rows.Scan(&p.ID, &p.Name, &p.Time); err != nil {
			return nil, err
		}

		// Simple logic for next prayer: first one with time > current time
		// This is a naive implementation assuming sorted order and today's times
		if !nextFound && p.Time > currentTime {
			p.IsNext = true
			nextFound = true
		}

		prayers = append(prayers, p)
	}

	// If no next prayer found (e.g. after Isya), next is Subuh tomorrow (first one)
	if !nextFound && len(prayers) > 0 {
		prayers[0].IsNext = true
	}

	return prayers, nil
}
