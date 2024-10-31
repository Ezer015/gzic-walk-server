package handlers

import (
	"gzic-walk-server/cache"
	"gzic-walk-server/config"
	"gzic-walk-server/database/db"
)

type Resolver struct {
	Conn   *db.Queries
	Caches *ResourceCaches
	Config *config.Configuration
}

type ResourceCaches struct {
	CopywritingCache *cache.TTLKeyedCache[string]
	ImageCache       *cache.TTLCache[int, []byte]
}
