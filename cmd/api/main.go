package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dmc0001/crud-project/internal/database"
	"github.com/dmc0001/crud-project/internal/env"
	"github.com/dmc0001/crud-project/internal/store"
	_ "github.com/lib/pq"
)

type Config struct {
	Addr         string
	newNoteModel *store.NoteModel
	errorLog     *log.Logger
	infoLog      *log.Logger
}
type Application struct {
	config Config
}

func main() {
	dsn := env.GetString("DB_DSN", "postgres://root:password@localhost:5432/name_db?sslmode=disable")
	port := env.GetString("PORT", ":3000")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := database.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cfg := Config{
		Addr:         port,
		newNoteModel: store.NewNoteModel(db),
		errorLog:     errorLog,
		infoLog:      infoLog,
	}
	app := &Application{
		config: cfg,
	}

	mux := http.NewServeMux()

	srv := app.route(mux)

	infoLog.Printf("Starting server on %s", app.config.Addr)

	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}

}
