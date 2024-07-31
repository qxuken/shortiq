package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/qxuken/short/internal"
	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/routes"
	"github.com/qxuken/short/web"
)

func main() {
	godotenv.Load()
	conf := internal.LoadConfig()
	if conf.Debug {
		log.Println("Application running in development mode")
		log.Printf("Config: %+v\n", conf)
	}

	dbPath := path.Join(conf.DataPath, "main.db?mode=rwc")
	log.Printf("Opeening db on %v\n", dbPath)
	db, err := db.ConnectSqlite3(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Recoverer)

	r.Get("/u/{short}", routes.RedirectRoute(db))
	r.Route("/a/e", routes.ExportRouter(db))
	r.Route("/api", routes.ApiRouter(conf, db))
	r.Group(web.WebRouter(conf, db))

	bind := fmt.Sprintf("%v:%v", conf.Bind, conf.Port)
	log.Printf("Listening on http://%v\n", bind)
	log.Printf("Available at %v\n", conf.PublicUrlStr)

	log.Fatal(http.ListenAndServe(bind, r))
}
