package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/thorbenbender/chirpy/internal/database"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	db, err := database.NewDB("./data/database.json")
	if err != nil {
		log.Fatal(err)
	}

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaApiKey := os.Getenv("POLKA_API_KEY")
	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
		JWTSecret:      jwtSecret,
		ApiKey:         polkaApiKey,
	}
	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))),
	)
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.HandleFunc("/reset", apiCfg.handleReset)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{id}", apiCfg.handlerChirpRetrieve)
	apiRouter.Delete("/chirps/{id}", apiCfg.handlerChirpDelete)
	apiRouter.Post("/users", apiCfg.handleUserCreate)
	apiRouter.Post("/login", apiCfg.handleUserLogin)
	apiRouter.Put("/users", apiCfg.handlerUserUpdate)
	apiRouter.Post("/refresh", apiCfg.HandleTokenRefresh)
	apiRouter.Post("/revoke", apiCfg.HandleTokenRevoke)
	apiRouter.Post("/polka/webhooks", apiCfg.HandlePolkaWebhook)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handleMetrics)

	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)
	corsMux := middlewareCors(router)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
