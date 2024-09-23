package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Suryarpan/rss-agg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// env setup
	godotenv.Load(".env")
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT not found in environment, exiting...")
	}
	log.Printf("Server is using port: %s\n", port)

	dbString, ok := os.LookupEnv("DB_URL")
	if !ok {
		log.Fatal("DB_URL not found in environment, exiting...")
	}
	log.Printf("Connected to: %s\n", dbString)

	// setup db connection
	conn, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("Could not connect to DB: %s\n", err)
	}
	queries := database.New(conn)
	apiCfg := apiConfig{
		DB: queries,
	}

	// router setup
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// nested router
	v1Router := chi.NewRouter()
	v1Router.Get("/health", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	
	// mount to original router
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("err")
	}
}
