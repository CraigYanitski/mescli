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
    secret     string
}

func main() {
    // load config file
    home := "."
    viper.SetConfigFile(path.Join(home, ".env"))
    viper.ReadInConfig()

    // get environment variables
    dbURL := viper.GetString("DB_URL")
    secret := viper.GetString("JWT_SECRET")
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
        secret: secret,
    }

    // create server multiplexer
    mux := http.NewServeMux()

    // define endpoint handlers
    // users
    mux.HandleFunc("POST /api/users", http.HandlerFunc(apiCfg.handleCreateUser))
    mux.HandleFunc("GET /api/users", http.HandlerFunc(apiCfg.handleGetUserByEmail))  //TODO
    mux.HandleFunc("PUT /api/users", http.HandlerFunc(apiCfg.handleUpdateUser))
    mux.HandleFunc("GET /api/users/{userID}", http.HandlerFunc(apiCfg.handleGetUser))  //TODO
    mux.HandleFunc("GET /api/users/crypto/{userID}", http.HandlerFunc(apiCfg.handleGetUserKeyPacket))
    // refresh tokens
    mux.HandleFunc("POST /api/login", http.HandlerFunc(apiCfg.handleLogin))
    mux.HandleFunc("POST /api/refresh", http.HandlerFunc(apiCfg.handleRefresh))
    mux.HandleFunc("POST /api/revoke", http.HandlerFunc(apiCfg.handleRevoke))
    // messages
    mux.HandleFunc("POST /api/messages", http.HandlerFunc(apiCfg.handleCreateMessage))
    mux.HandleFunc("GET /api/messages", http.HandlerFunc(apiCfg.HandleGetMessages))

    // define server and listen for requests
    const port = "8080"
    server := http.Server{
        Addr: ":" + port,
        Handler: mux,
    }
    fmt.Printf("Serving mescli api on port: %v\n", port)
    log.Fatal(server.ListenAndServe())
}

