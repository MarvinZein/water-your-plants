package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Plant struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Place         string    `json:"place"`
	LastWateredAt time.Time `json:"lastWateredAt"`
}

func main() {
	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create plants table if not exists
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS plants (id SERIAL PRIMARY KEY, name TEXT, place TEXT, last_watered_at DATE);")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/go/plants", getPlants(db)).Methods("GET")
	router.HandleFunc("/api/go/plants", createPlant(db)).Methods("POST")
	router.HandleFunc("/api/go/plants/{id}", getPlant(db)).Methods("GET")
	router.HandleFunc("/api/go/plants/{id}", updatePlant(db)).Methods("PUT")
	router.HandleFunc("/api/go/plants/{id}", deletePlant(db)).Methods("DELETE")

	// wrap router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Get all plants
func getPlants(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM plants;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		plants := []Plant{}
		for rows.Next() {
			var p Plant
			if err := rows.Scan(&p.Id, &p.Name, &p.Place, &p.LastWateredAt); err != nil {
				log.Fatal(err)
			}
			plants = append(plants, p)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		response := map[string][]Plant{"plants": plants}
		json.NewEncoder(w).Encode(response)
	}
}

// Get plant by id
func getPlant(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var p Plant
		err := db.QueryRow("SELECT * FROM plants WHERE id = $1", id).Scan(&p.Id, &p.Name, &p.Place, &p.LastWateredAt)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response := map[string]Plant{"plant": p}
		json.NewEncoder(w).Encode(response)
	}
}

// Creates a new plant
func createPlant(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p Plant
		json.NewDecoder(r.Body).Decode(&p)

		err := db.QueryRow("INSERT INTO plants (name, place, last_watered_at) VALUES ($1, $2, $3) RETURNING id", p.Name, p.Place, p.LastWateredAt).Scan(&p.Id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(p)
	}
}

// Updates an existing plant
func updatePlant(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p Plant
		json.NewDecoder(r.Body).Decode(&p)

		vars := mux.Vars(r)
		id := vars["id"]

		// Execute the update query
		_, err := db.Exec("UPDATE plants SET name = $1, place = $2, last_watered_at = $3 WHERE id = $4", p.Name, p.Place, p.LastWateredAt, id)
		if err != nil {
			log.Fatal(err)
		}

		// Retrieve the updated plant data from the database
		var updatedPlant Plant
		err = db.QueryRow("SELECT * FROM plants WHERE id = $1", id).Scan(&updatedPlant.Id, &updatedPlant.Name, &updatedPlant.Place, &updatedPlant.LastWateredAt)
		if err != nil {
			log.Fatal(err)
		}

		// Send the updated plant data in the response
		json.NewEncoder(w).Encode(updatedPlant)
	}
}

// delete plant
func deletePlant(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var p Plant
		err := db.QueryRow("SELECT * FROM plants WHERE id = $1", id).Scan(&p.Id, &p.Name, &p.Place, &p.LastWateredAt)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM plants WHERE id = $1", id)
			if err != nil {
				//todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Plant deleted")
		}
	}
}
