package storage

import (
	"log"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/Slayzur02/GoChess/pkg/internal_sql"
	redis "github.com/go-redis/redis/v8"
)

type Repo struct {
	mc *redis.Client
	db *internal_sql.Queries
}

// instantiates redis
func instantiateRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

// instantiates database with queries
func instantiateDB() *internal_sql.Queries {
	db, err := sql.Open("postgres", "user=macbookpro dbname=eski_chess sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	return internal_sql.New(db)
}

// return new repo
func NewRepo() *Repo {
	return &Repo{
		mc: instantiateRedis(),
		db: instantiateDB(),
	}
}
