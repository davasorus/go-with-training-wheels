package main

import (
	"fmt"
	"log"
	"log/slog"

	database "example/data-access/Database"
)

func main() {

	database.ConnectDatabase()

	albums, err := database.AlbumsByArtist("John Coltrane")
	if err != nil {
		slog.Error("Error fetching albums by artist", "error", err)
		return
	}

	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := database.AlbumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := database.AddAlbum(database.Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}
