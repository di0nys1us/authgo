package authgo

import (
	"github.com/jmoiron/sqlx"
)

type db struct {
	*sqlx.DB
}

type tx struct {
	*sqlx.Tx
}
