package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gzic-walk-server/database/db"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Resolver) CreateRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateRecord called")
	// Parse the form data
	imageID, err := strconv.Atoi(r.FormValue("image_id"))
	if err != nil {
		log.Println("Invalid image ID:", err)
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}
	copywriting := r.FormValue("copywriting")
	if copywriting == "" {
		log.Println("Copywriting cannot be empty")
		http.Error(w, "Copywriting cannot be empty", http.StatusBadRequest)
		return
	}
	var sightID int
	sightIDStr := r.FormValue("sight_id")
	sightName := r.FormValue("sight_name")
	switch {
	case sightIDStr == "" && sightName == "":
		log.Println("Sight ID or sight name must be provided")
		http.Error(w, "Sight ID or sight name must be provided", http.StatusBadRequest)
		return
	case sightIDStr != "" && sightName != "":
		log.Println("Sight ID and sight name cannot both be provided")
		http.Error(w, "Sight ID and sight name cannot both be provided", http.StatusBadRequest)
		return
	case sightIDStr != "":
		sightID, err = strconv.Atoi(sightIDStr)
		if err != nil {
			log.Println("Invalid sight ID:", err)
			http.Error(w, "Invalid sight ID", http.StatusBadRequest)
			return
		}
		// Validate the sight id
		sights, err := s.Conn.GetSights(r.Context())
		if err != nil {
			log.Println("Failed to validate sight id:", err)
			http.Error(w, "Failed to validate sight id", http.StatusInternalServerError)
			return
		}
		valid := false
		for _, sight := range sights {
			// Check if the sight name already exists
			if int32(sightID) == sight.SightID {
				valid = true
				break
			}
		}
		if !valid {
			log.Println("Invalid sight id")
			http.Error(w, "Invalid sight id", http.StatusBadRequest)
			return
		}
	case sightName != "":
		// Validate the sight name
		sights, err := s.Conn.GetSights(r.Context())
		if err != nil {
			log.Println("Failed to validate sight name:", err)
			http.Error(w, "Failed to validate sight name", http.StatusInternalServerError)
			return
		}
		for _, sight := range sights {
			// Check if the sight name already exists
			if sightName == sight.SightName {
				log.Println("Invalid sight name")
				http.Error(w, "Invalid sight name", http.StatusBadRequest)
				return
			}
		}
	default:
		log.Println("Switch case did not match any valid condition")
		http.Error(w, "Switch case did not match any valid condition", http.StatusBadRequest)
	}

	// Create the record
	reqBody := db.CreateRecordParams{
		ImageID: int32(imageID),
		SightID: pgtype.Int4{
			Int32: int32(sightID),
			Valid: sightIDStr != "",
		},
		SightName:   sightName,
		Copywriting: copywriting,
	}
	recordID, err := s.Conn.CreateRecord(r.Context(), reqBody)
	if err != nil {
		log.Println("Failed to create record:", err)
		http.Error(w, "Failed to create record", http.StatusInternalServerError)
		return
	}

	// Return the record ID as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Println("CreateRecord successful")
	json.NewEncoder(w).Encode(map[string]any{
		"record_id": recordID,
	})
}

func (s *Resolver) GetRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("GetRecord called")
	// Get the record ID from the URL
	recordID, err := strconv.Atoi(r.PathValue("record_id"))
	if err != nil {
		log.Println("Invalid record ID:", err)
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	// Get the record from the database
	record, err := s.Conn.GetRecord(r.Context(), int32(recordID))
	if err != nil {
		log.Println("Record not found:", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Return the record as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("GetRecord successful")
	json.NewEncoder(w).Encode(record)
}

func (s *Resolver) GetRandomRecord(w http.ResponseWriter, r *http.Request) {
	log.Println("GetRandomRecord called")
	// Get a random record from the database
	record, err := s.Conn.GetRandomRecord(r.Context())
	if err != nil {
		log.Println("Failed to get random record:", err)
		http.Error(w, "Failed to get random record", http.StatusInternalServerError)
		return
	}

	// Return the record as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("GetRandomRecord successful")
	json.NewEncoder(w).Encode(record)
}
