package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"

	"gzic-walk-server/cache"
	"gzic-walk-server/config"
	"gzic-walk-server/database/db"
	"gzic-walk-server/handlers"
)

func main() {
	// Set log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Load the configuration
	configFilePath := os.Getenv("CONFIG_PATH")
	if configFilePath == "" {
		log.Fatalln("CONFIG_PATH environment variable is required")
	}
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalln("Unable to load configuration: ", err)
	}

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatalln("Unable to connect to database: ", err)
	}
	defer conn.Close(context.Background())

	// Initialize the resolver
	s := &handlers.Resolver{
		Conn: db.New(conn),
		Caches: &handlers.ResourceCaches{
			CopywritingCache: cache.NewTTLKeyedCache[string](2 * handlers.CopywritingExpiry),
			ImageCache:       cache.NewTTLCache[int, []byte](2 * handlers.ImageExpiry),
		},
		Config: config,
	}

	// Register the handlers
	http.HandleFunc("POST /image", s.UploadImage)
	http.HandleFunc("GET /image/{image_id}", s.DownloadImage)
	http.HandleFunc("GET /sight", s.GetSights)
	http.HandleFunc("GET /sight/{sight_id}", s.GetSight)
	http.HandleFunc("POST /copywriting", s.CreateCopywriting)
	http.HandleFunc("GET /copywriting/{copywriting_id}", s.GetCopywriting)
	http.HandleFunc("POST /record", s.CreateRecord)
	http.HandleFunc("GET /record/{record_id}", s.GetRecord)
	http.HandleFunc("GET /record", s.GetRandomRecord)

	// Start the server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Unable to start server: ", err)
	}
}
