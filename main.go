package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thorbenbender/chirpy/internal/database"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	db, err := database.NewDB("./data/database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
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
	apiRouter.Post("/users", apiCfg.handleUserCreate)
	apiRouter.Post("/login", apiCfg.handleUserLogin)

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
