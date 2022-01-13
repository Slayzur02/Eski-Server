package storage

import (
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
func instantiateDB(db internal_sql.DBTX) *internal_sql.Queries {
	return internal_sql.New(db)
}

// return new repo
func NewRepo(db internal_sql.DBTX) *Repo {
	return &Repo{
		mc: instantiateRedis(),
		db: instantiateDB(db),
	}
}
