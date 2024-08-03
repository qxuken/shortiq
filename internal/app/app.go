package app

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/qxuken/short/internal/config"
	"github.com/qxuken/short/internal/db"
)

func RunApp() {
	conf := config.LoadConfig()

	dbPath := path.Join(conf.DataPath, "main.db?mode=rwc")
	db := db.ConnectSqlite3(conf, dbPath)

	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Recoverer)

	appRouter(conf, db, r)

	bind := fmt.Sprintf("%v:%v", conf.Bind, conf.Port)
	log.Printf("Listening on http://%v\n", bind)
	log.Printf("Available at %v\n", conf.PublicUrlStr)

	log.Fatal(http.ListenAndServe(bind, r))
}
