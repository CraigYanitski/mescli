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
    mux.Handle("GET /api/users", apiCfg.authenticationMiddleware(http.HandlerFunc(apiCfg.handleGetUserByEmail)))
    mux.Handle("PUT /api/users", apiCfg.authenticationMiddleware(http.HandlerFunc(apiCfg.handleUpdateUser)))
    mux.Handle("GET /api/users/{userID}", apiCfg.authenticationMiddleware(http.HandlerFunc(apiCfg.handleGetUser)))
    mux.Handle("GET /api/users/crypto/{userID}", apiCfg.authenticationMiddleware(http.HandlerFunc(apiCfg.handleGetUserKeyPacket)))
    // refresh tokens
    mux.HandleFunc("POST /api/login", http.HandlerFunc(apiCfg.handleLogin))
    mux.HandleFunc("POST /api/refresh", http.HandlerFunc(apiCfg.handleRefresh))
    mux.HandleFunc("POST /api/revoke", http.HandlerFunc(apiCfg.handleRevoke))
    // messages
    mux.Handle("POST /api/messages", apiCfg.authenticationMiddleware(http.HandlerFunc(apiCfg.handleCreateMessage)))
    mux.Handle("GET /api/messages", apiCfg.authenticationMiddleware(http.HandlerFunc(apiCfg.HandleGetMessages)))

    // define server and listen for requests
    const port = "8080"
    server := http.Server{
        Addr: ":" + port,
        Handler: mux,
    }
    fmt.Printf("Serving mescli api on port: %v\n", port)
    log.Fatal(server.ListenAndServe())
}

