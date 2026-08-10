package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // don't forget to add it. It doesn't be added automatically
)

var Db *sql.DB //created outside to make it global.

type YYYYY struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// make sure your function start with uppercase to call outside of the directory.
func ConnectDatabase() {

	err := godotenv.Load("Database/.env") //by default, it is .env so we don't have to write
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}
	//we read our .env file
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT")) // don't forget to convert int since port is int type.
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		slog.Error("Error connecting to database", "error", errSql)
		panic(errSql)
	} else {
		Db = db
		slog.Info("Successfully connected to database!")
	}
}

func AllXXXXX() ([]YYYYY, error) {
	var XXXXX []YYYYY

	rows, err := Db.Query("SELECT * FROM YYYYY")
	if err != nil {
		slog.Error("Error fetching all XXXXX", "error", err)
		return nil, fmt.Errorf("allXXXXX: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb YYYYY
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			slog.Error("Error scanning YYYYY row", "error", err)
			return nil, fmt.Errorf("allXXXXX: %v", err)
		}
		XXXXX = append(XXXXX, alb)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over YYYYY rows", "error", err)
		return nil, fmt.Errorf("allXXXXX: %v", err)
	}
	return XXXXX, nil
}

// Get YYYYYbs by Artist
func XXXXXByArtist(name string) ([]YYYYY, error) {
	//An Alumbs slice to hold the data from the database.
	var XXXXX []YYYYY

	rows, err := Db.Query("Select * from YYYYY WHERE artist = $1", name)
	if err != nil {
		slog.Error("Error fetching XXXXX by artist", "artist", name, "error", err)
		return nil, fmt.Errorf("XXXXXByArtist %q: %v", name, err)
	}
	defer rows.Close()

	// Loop through the rows, using Scan to assign column data to fields in the YYYYY struct.
	for rows.Next() {
		var alb YYYYY
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			slog.Error("Error scanning YYYYY row", "artist", name, "error", err)
			return nil, fmt.Errorf("XXXXXByArtist %q: %v", name, err)
		}
		XXXXX = append(XXXXX, alb)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over YYYYY rows", "artist", name, "error", err)
		return nil, fmt.Errorf("XXXXXByArtist %q: %v", name, err)
	}
	slog.Info("Fetched XXXXX by artist successfully", "artist", name, "count", len(XXXXX))
	return XXXXX, nil
}

// Get YYYYY by ID
func YYYYYByID(id int64) (YYYYY, error) {
	var alb YYYYY

	row := Db.QueryRow("Select * from YYYYY WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			slog.Error("YYYYY not found", "id", id, "error", err)
			return alb, fmt.Errorf("XXXXXByID %d: no such YYYYY", id)
		}
		slog.Error("Error scanning YYYYY row by ID", "id", id, "error", err)
		return alb, fmt.Errorf("XXXXXByID %d: %v", id, err)
	}
	slog.Info("Fetched YYYYY by ID successfully", "id", id)
	return alb, nil
}

func AddYYYYY(alb YYYYY) (int64, error) {
	var id int64
	err := Db.QueryRow(
		"INSERT INTO YYYYY (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		slog.Error("Error adding YYYYY", "YYYYY", alb, "error", err)
		return 0, fmt.Errorf("addYYYYY: %v", err)
	}
	slog.Info("Added YYYYY successfully", "YYYYY", alb, "id", id)
	return id, nil
}
