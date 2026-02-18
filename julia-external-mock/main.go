package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// NewsExternalItem matches the model in julia-notification-batch
type NewsExternalItem struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Channel     string `json:"channel"`
}

func main() {
	http.HandleFunc("/api/v1/external/news", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)

		news := []NewsExternalItem{
			{
				ID:          "EXT-1",
				Description: "Nuova migrazione Julia completata con successo!",
				Channel:     "INTERNAL",
			},
			{
				ID:          "EXT-2",
				Description: "Aggiornamento disponibile per il modulo batch.",
				Channel:     "EXTERNAL",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(news)
	})

	port := ":8096"
	fmt.Printf("Mock External News Provider starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
