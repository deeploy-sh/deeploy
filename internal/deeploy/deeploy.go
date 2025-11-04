package deeploy

import (
	"database/sql"
	"net/http"
)

type App struct {
	Router *http.ServeMux
	DB     *sql.DB
}

func New(router *http.ServeMux, db *sql.DB) App {
	return App{Router: router, DB: db}
}
