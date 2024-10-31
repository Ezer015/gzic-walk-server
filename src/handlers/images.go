package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	maxFileSize = 4 << 20 // 4 MB
	ImageExpiry = time.Minute
)

func (s *Resolver) UploadImage(w http.ResponseWriter, r *http.Request) {
	log.Println("UploadImage called")
	// Retrieve the image from the form data
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	file, _, err := r.FormFile("image")
	if err != nil {
		if err.Error() == "http: request body too large" {
			log.Println("Image exceeds maximum file size")
			http.Error(w, "Image exceeds maximum file size", http.StatusRequestEntityTooLarge)
			return
		}
		log.Println("Error retrieving the image:", err)
		http.Error(w, "Error retrieving the image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the storage directory if it doesn't exist
	if _, err := os.Stat(s.Config.StoragePath); os.IsNotExist(err) {
		err = os.MkdirAll(s.Config.StoragePath, os.ModePerm)
		if err != nil {
			log.Println("Failed to create storage directory:", err)
			http.Error(w, "Failed to create storage directory", http.StatusInternalServerError)
			return
		}
	}

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("Failed to decode image:", err)
		http.Error(w, "Failed to decode image", http.StatusBadRequest)
		return
	}
	// Generate JPEG file name
	fileName := fmt.Sprintf("%d_%s.jpg", time.Now().Unix(), uuid.New().String())
	dst, err := os.Create(filepath.Join(s.Config.StoragePath, fileName))
	if err != nil {
		log.Println("Failed to create file:", err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	// Encode the image as JPEG
	if err := jpeg.Encode(dst, img, &jpeg.Options{Quality: 75}); err != nil {
		os.Remove(dst.Name())
		log.Println("Failed to save image as JPEG:", err)
		http.Error(w, "Failed to save image as JPEG", http.StatusInternalServerError)
		return
	}

	// Save the image path to the database
	_, err = s.Conn.CreateImage(r.Context(), pgtype.Text{String: fileName, Valid: true})
	if err != nil {
		os.Remove(dst.Name())
		log.Println("Failed to save image:", err)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Register the image for processing
	processingID, err := s.Conn.CreateImage(r.Context(), pgtype.Text{})
	if err != nil {
		log.Println("Failed to prepare image processing:", err)
		http.Error(w, "Failed to prepare image processing", http.StatusInternalServerError)
		return
	}

	// Process the image in the background
	go func() {
		log.Println("Image processing started for ID:", processingID)
		log.Println("Image processing completed for ID:", processingID)
	}()

	// Return the image ID as JSON
	w.WriteHeader(http.StatusAccepted)
	log.Println("UploadImage successful")
	json.NewEncoder(w).Encode(map[string]any{
		"image_id": processingID,
	})
}

func (s *Resolver) DownloadImage(w http.ResponseWriter, r *http.Request) {
	log.Println("DownloadImage called")
	// Get the image ID from the URL
	imageID, err := strconv.Atoi(r.PathValue("image_id"))
	if err != nil {
		log.Println("Invalid image ID:", err)
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Check the image cache
	imageData, found := s.Caches.ImageCache.Get(imageID)

	if !found {
		// Get the image path from the database
		imagePath, err := s.Conn.GetImage(r.Context(), int32(imageID))
		if err != nil {
			log.Println("Image not found:", err)
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		if !imagePath.Valid {
			log.Println("Image not ready")
			http.Error(w, "Image not ready", http.StatusAccepted)
			return
		}

		// Open the image file
		file, err := os.Open(filepath.Join(s.Config.StoragePath, imagePath.String))
		if err != nil {
			log.Println("Failed to open image:", err)
			http.Error(w, "Failed to open image", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Read the image data
		imageData, err = io.ReadAll(file)
		if err != nil {
			log.Println("Failed to read image:", err)
			http.Error(w, "Failed to read image", http.StatusInternalServerError)
			return
		}
		// Cache the image data
		s.Caches.ImageCache.Set(imageID, imageData, ImageExpiry)
	}

	// Return the image
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	log.Println("DownloadImage successful")
	w.Write(imageData)
}
