package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/qxuken/short/internal"
	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/web"
)

func main() {
	godotenv.Load()
	conf := internal.LoadConfig()
	if conf.Debug {
		log.Println("Application runnin in development mode")
	}

	db, err := db.ConnectSqlite3("./tmp/db.db?mode=rwc")
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Recoverer)

	r.Get("/u/{short}", internal.RedirectRoute(db))
	r.Group(web.WebRouter(db))

	log.Printf("Listening on http://127.0.0.1:%v\n", conf.Port)

	log.Fatal(http.ListenAndServe(":8080", r))
}
