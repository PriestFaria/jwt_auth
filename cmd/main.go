package main

import (
	"auth_backend/infrastructure"
	"auth_backend/infrastructure/dataAccess"
	"auth_backend/infrastructure/middleware"
	"auth_backend/presentation"
	"auth_backend/usecase"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	databaseConnectionString := os.Getenv("DATABASE_CONNECTION_STRING")
	if databaseConnectionString == "" {
		databaseConnectionString = os.Getenv("DATABASE_CONNECTION_STRING")
	}

	db, err := sql.Open("postgres", databaseConnectionString)
	if err != nil {
		log.Fatal("error connecting to the database %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("error connecting to the database %v\n", err)
	}
	userRepository := dataAccess.NewUserRepositoryDB(db)

	tokenGenerator := infrastructure.NewJwtTokenGenerator("secret")

	userUseCase := usecase.NewUserUseCase(userRepository, tokenGenerator)

	userHandler := presentation.NewUserHandler(userUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/register", userHandler.Register)

	adminHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := middleware.GetClaims(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Добро пожаловать, администратор!", "email": "` + claims["email"].(string) + `"}`))
	})

	secret := "secret"
	handlerChain := middleware.JwtMiddleware(secret)(middleware.AdminMiddleware(adminHandler))
	mux.Handle("/admin", handlerChain)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
