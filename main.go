package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/CraigYanitski/mescli/internal/database"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type apiConfig struct {
    dbQueries  *database.Queries
}

func main() {
    // load config file
    home := "."
    viper.SetConfigFile(path.Join(home, ".env"))
    viper.ReadInConfig()

    // get environment variables
    dbURL := viper.GetString("DB_URL")
    log.Println(dbURL)

    // open database
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("cannot open database %s: %s", dbURL, err)
    }
    dbQueries := database.New(db)

    // define database persistent configuration
    apiCfg := apiConfig{
        dbQueries: dbQueries,
    }

    // create server multiplexer
    mux := http.NewServeMux()

    // define endpoint handlers
    mux.HandleFunc("POST /api/users", http.HandlerFunc(apiCfg.handleCreateUser))
    mux.HandleFunc("GET /api/users", http.HandlerFunc(apiCfg.handleGetUserKeyPacket))

    // define server and listen for requests
    const port = "8080"
    server := http.Server{
        Addr: ":" + port,
        Handler: mux,
    }
    fmt.Printf("Serving mescli api on port: %v\n", port)
    log.Fatal(server.ListenAndServe())
}

