// Code generated by sqlc. DO NOT EDIT.

package internal_sql

import (
	"database/sql"
)

type AppUser struct {
	ID        int64
	Email     string
	Username  string
	Hashedpwd string
}

type Game struct {
	ID      int64
	Whiteid sql.NullInt64
	Blackid sql.NullInt64
}

type GameMove struct {
	TurnNumber     sql.NullInt32
	Color          sql.NullInt32
	Startsquare    sql.NullString
	Endsquare      sql.NullString
	Promotionpiece sql.NullString
	GameID         sql.NullInt64
}
