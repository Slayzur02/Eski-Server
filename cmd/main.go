package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"

	publicHttp "net/http"

	"github.com/Slayzur02/GoChess/pkg/auth"
	"github.com/Slayzur02/GoChess/pkg/http"
	"github.com/Slayzur02/GoChess/pkg/storage"
)

func main() {
	// load environment and get keys
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	user := os.Getenv("user")

	// get db
	db, err := sql.Open("postgres", "user="+user+" dbname=eski_chess sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := storage.NewRepo(db)

	a := auth.NewService(repo)
	router := http.HandleCreation(a)

	log.Fatal(publicHttp.ListenAndServe(":8080", router))
}
