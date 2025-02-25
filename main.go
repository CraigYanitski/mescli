package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CraigYanitski/mescli/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
    dbQueries  *database.Queries
}

func main() {
    godotenv.Load()
    dbURL := os.Getenv("DB_URL")

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("cannot open database %s: %s", dbURL, err)
    }
    dbQueries := database.New(db)

    apiCfg := apiConfig{
        dbQueries: dbQueries,
    }

    mux := http.NewServeMux()

    // Handlers
    mux.HandleFunc("POST /api/users", http.HandlerFunc(apiCfg.handleCreateUser))

    //const fsPATH = "."
    const port = "8080"
    server := http.Server{
        Addr: ":" + port,
        Handler: mux,
    }

    fmt.Printf("Serving mescli api on port: %v\n", port)
    log.Fatal(server.ListenAndServe())
}

