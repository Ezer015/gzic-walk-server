package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (s *Resolver) GetSights(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSights called")
	// Get the sights from the database
	sights, err := s.Conn.GetSights(r.Context())
	if err != nil {
		log.Println("Failed to get sights:", err)
		http.Error(w, "Failed to get sights", http.StatusInternalServerError)
		return
	}

	// Return the sights as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("GetSights successful")
	json.NewEncoder(w).Encode(sights)
}

func (s *Resolver) GetSight(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSight called")
	// Get the sight ID from the URL
	sightID, err := strconv.Atoi(r.PathValue("sight_id"))
	if err != nil {
		log.Println("Invalid sight ID:", err)
		http.Error(w, "Invalid sight ID", http.StatusBadRequest)
		return
	}

	// Get the sight from the database
	sight, err := s.Conn.GetSight(r.Context(), int32(sightID))
	if err != nil {
		log.Println("Sight not found:", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Return the sight as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("GetSight successful")
	json.NewEncoder(w).Encode(sight)
}
