package main

import (
	"log"

	publicHttp "net/http"

	"github.com/Slayzur02/GoChess/pkg/auth"
	"github.com/Slayzur02/GoChess/pkg/http"
	"github.com/Slayzur02/GoChess/pkg/storage"
)

func main() {
	repo := storage.NewRepo()
	a := auth.NewService(repo)

	router := http.HandleCreation(&a)

	log.Fatal(publicHttp.ListenAndServe(":8080", router))
}
