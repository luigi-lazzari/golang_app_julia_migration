package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// NewsItem matches the model in julia-notification-batch
type NewsItem struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Channel     string `json:"channel"`
}

func main() {
	http.HandleFunc("/api/v1/news", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)

		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var news []NewsItem
		if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Printf("Successfully received %d news items:", len(news))
		for _, item := range news {
			fmt.Printf("-----------------------------------\n"+
				"ID:          %s\n"+
				"Channel:     %s\n"+
				"Description: %s\n"+
				"-----------------------------------\n",
				item.ID, item.Channel, item.Description)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "News updated successfully")
	})

	port := ":8095"
	fmt.Printf("Mock Internal Notification Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
