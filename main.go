package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Demianeen/chirpy/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	const filepathRoot = "./static"
	const dbPath = "./database.json"
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		os.Remove(dbPath)
	}

	newDb, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	config := apiConfig{
		fileserverHits: 0,
		db:             newDb,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		polkaApiKey:    os.Getenv("POLKA_API_KEY"),
	}
	mux := http.NewServeMux()

	mux.Handle("GET /app/*", http.StripPrefix("/app/", config.middlewareMetrics(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./static/admin/metrics/index.html")
		if err != nil {
			log.Println("Error reading file:", err)
			return
		}
		htmlString := string(data)

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, htmlString, config.fileserverHits)
	})

	mux.HandleFunc("POST /api/chirps", config.handleCreateChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", config.handlerDeleteChirp)
	mux.HandleFunc("GET /api/chirps", config.handleRetrieveChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", config.handleGetChirp)

	mux.HandleFunc("POST /api/users", config.handleCreateUser)
	mux.HandleFunc("PUT /api/users", config.handleUpdateUser)
	mux.HandleFunc("POST /api/login", config.handleLoginUser)

	mux.HandleFunc("POST /api/refresh", config.handleRefreshJwtToken)
	mux.HandleFunc("POST /api/revoke", config.handleRevokeRefreshToken)

	mux.HandleFunc("POST /api/polka/webhooks", config.handlePolkaWebhook)

	mux.HandleFunc("GET /api/reset", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		config.fileserverHits = 0
		w.Write([]byte("Hits reset to 0"))
	})

	corsMux := middlewareCors(mux)
	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
