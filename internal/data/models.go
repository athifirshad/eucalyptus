package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	Users UserModel
	Tokens TokenModel
}

func NewModels(dbPool *pgxpool.Pool) Models {
	return Models{
		Users: UserModel{dbPool: dbPool},
		Tokens: TokenModel{dbPool: dbPool},
	}
}
