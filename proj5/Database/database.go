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

type Album struct {
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

func AllAlbums() ([]Album, error) {
	var albums []Album

	rows, err := Db.Query("SELECT * FROM album")
	if err != nil {
		slog.Error("Error fetching all albums", "error", err)
		return nil, fmt.Errorf("allAlbums: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			slog.Error("Error scanning album row", "error", err)
			return nil, fmt.Errorf("allAlbums: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over album rows", "error", err)
		return nil, fmt.Errorf("allAlbums: %v", err)
	}
	return albums, nil
}

// Get Albumbs by Artist
func AlbumsByArtist(name string) ([]Album, error) {
	//An Alumbs slice to hold the data from the database.
	var albums []Album

	rows, err := Db.Query("Select * from album WHERE artist = $1", name)
	if err != nil {
		slog.Error("Error fetching albums by artist", "artist", name, "error", err)
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	// Loop through the rows, using Scan to assign column data to fields in the Album struct.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			slog.Error("Error scanning album row", "artist", name, "error", err)
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over album rows", "artist", name, "error", err)
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	slog.Info("Fetched albums by artist successfully", "artist", name, "count", len(albums))
	return albums, nil
}

// Get Album by ID
func AlbumByID(id int64) (Album, error) {
	var alb Album

	row := Db.QueryRow("Select * from album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Album not found", "id", id, "error", err)
			return alb, fmt.Errorf("albumsByID %d: no such album", id)
		}
		slog.Error("Error scanning album row by ID", "id", id, "error", err)
		return alb, fmt.Errorf("albumsByID %d: %v", id, err)
	}
	slog.Info("Fetched album by ID successfully", "id", id)
	return alb, nil
}

func AddAlbum(alb Album) (int64, error) {
	var id int64
	err := Db.QueryRow(
		"INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		slog.Error("Error adding album", "album", alb, "error", err)
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	slog.Info("Added album successfully", "album", alb, "id", id)
	return id, nil
}
