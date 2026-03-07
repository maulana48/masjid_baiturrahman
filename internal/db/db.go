package db

import (
	"database/sql"
	"log"
	"masjid_baiturrahman/internal/models"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./masjid.db"
	}
	DB, err = sql.Open("sqlite3", dbPath)
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

	createContactsTable := `
	CREATE TABLE IF NOT EXISTS contacts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		phone TEXT,
		email TEXT,
		topic TEXT,
		message TEXT,
		created_at DATETIME
	);
	`
	_, err = DB.Exec(createContactsTable)
	if err != nil {
		log.Fatal(err)
	}

	createProgramsTable := `
	CREATE TABLE IF NOT EXISTS programs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		category TEXT,
		arabic_title TEXT,
		description TEXT,
		ustadz TEXT,
		schedule TEXT,
		level TEXT,
		quota TEXT,
		is_featured BOOLEAN DEFAULT 0,
		show_on_home BOOLEAN DEFAULT 0
	);
	`
	_, err = DB.Exec(createProgramsTable)
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

	// Seed programs if empty
	var progCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM programs").Scan(&progCount)
	if err != nil {
		log.Fatal(err)
	}

	if progCount == 0 {
		insertProgs := `
		INSERT INTO programs (title, category, arabic_title, description, ustadz, schedule, level, quota, is_featured, show_on_home) VALUES
		('Tahfizh Intensif', 'unggulan', 'تَحْفِيْظ', 'Program hafalan Al-Quran intensif dengan target 30 juz.', 'Ust. Qari Ridho', 'Setiap Hari', 'Semua Usia', 'Terbatas', 1, 0),
		('Kajian Tafsir', 'kajian', 'تَفْسِيْر', 'Membedah makna ayat-ayat Al-Quran secara mendalam.', 'Ust. Dr. Fauzan', 'Senin & Kamis Malam', 'Umum', 'Terbuka', 0, 1),
		('TPA Anak', 'pendidikan', 'تَعْلِيْم', 'Pendidikan Al-Quran dasar untuk anak-anak.', 'Tim TPA', 'Sabtu & Minggu Pagi', 'Anak-anak', '50 Santri', 0, 1),
		('Remaja Masjid', 'remaja', 'شَبَاب', 'Wadah kreativitas dan dakwah remaja.', 'Pembina Remaja', 'Jumat Sore', 'Remaja', 'Terbuka', 0, 0),
		('Baksos & Zakat', 'sosial', 'زَكَاة', 'Penyaluran bantuan sosial dan pengelolaan zakat.', 'Panitia ZIS', 'Kondisional', 'Umum', '-', 0, 1),
		('Fiqh Muamalah', 'kajian', 'فِقْه', 'Kajian hukum ekonomi dan transaksi Islam.', 'Ust. Syarifuddin', 'Selasa Pagi', 'Umum', 'Terbuka', 0, 0);
		`
		_, err = DB.Exec(insertProgs)
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

func GetAllPrograms() ([]models.Program, error) {
	rows, err := DB.Query("SELECT id, title, category, arabic_title, description, ustadz, schedule, level, quota, is_featured, show_on_home FROM programs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		if err := rows.Scan(&p.ID, &p.Title, &p.Category, &p.ArabicTitle, &p.Description, &p.Ustadz, &p.Schedule, &p.Level, &p.Quota, &p.IsFeatured, &p.ShowOnHome); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, nil
}

func GetHomePrograms() ([]models.Program, error) {
	rows, err := DB.Query("SELECT id, title, category, arabic_title, description, ustadz, schedule, level, quota, is_featured, show_on_home FROM programs WHERE show_on_home = 1 LIMIT 3")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		if err := rows.Scan(&p.ID, &p.Title, &p.Category, &p.ArabicTitle, &p.Description, &p.Ustadz, &p.Schedule, &p.Level, &p.Quota, &p.IsFeatured, &p.ShowOnHome); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, nil
}

func SaveContact(msg models.ContactMessage) error {
	stmt, err := DB.Prepare("INSERT INTO contacts(name, phone, email, topic, message, created_at) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(msg.Name, msg.Phone, msg.Email, msg.Topic, msg.Message, msg.CreatedAt)
	return err
}
