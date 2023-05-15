package model

import "database/sql"

type User struct {
	Id              int64
	Name            string         `json:"name"`
	Username        string         `json:"username"`
	Password        string         `json:"password"`
	RefreshToken    sql.NullString `json:"refresh_token"`
	RefreshTokenEAT sql.NullInt64  `json:"refresh_token_eat"`
	Role            string         `json:"role"`
}
