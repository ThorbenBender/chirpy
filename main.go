package main

import (
	"log"
	"net/http"
  "github.com/go-chi/chi/v5"
)

func main() {
  const filePathRoot = "."
  const port = "8080"
  apiCfg := apiConfig{
    fileServerHits: 0,
  } 
  router := chi.NewRouter()
  fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
  router.Handle("/app/*", fsHandler)
  router.Handle("/app", fsHandler)
  
  apiRouter := chi.NewRouter()
  apiRouter.Get("/healthz", handleReadiness)
  apiRouter.HandleFunc("/reset", apiCfg.handleReset)
  apiRouter.Post("/validate_chirp", validate_chirp)

  adminRouter := chi.NewRouter()
  adminRouter.Get("/metrics", apiCfg.handleMetrics)



  router.Mount("/api", apiRouter)
  router.Mount("/admin", adminRouter)
  corsMux := middlewareCors(router)
  server := &http.Server{
    Addr: ":" + port,
    Handler: corsMux,
  }
  log.Printf("Serving on port: %s\n", port)
  log.Fatal(server.ListenAndServe())
}
