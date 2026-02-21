package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port            string
	MongoURI        string
	MongoDB         string
	JWTSecret       string
	JWTAccessTTL    int
	JWTRefreshTTL   int
	GCSBucket       string
	GCPProjectID    string
	VisionLocation  string
	YouTubeAPIKey   string
	GoogleSearchKey string
	GoogleSearchCX  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		MongoURI:        getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGODB_DB", "fridge_master"),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		GCSBucket:       getEnv("GCS_BUCKET", ""),
		GCPProjectID:    getEnv("GCP_PROJECT_ID", ""),
		VisionLocation:  getEnv("VISION_LOCATION", "us-central1"),
		YouTubeAPIKey:   getEnv("YOUTUBE_API_KEY", ""),
		GoogleSearchKey: getEnv("GOOGLE_SEARCH_KEY", ""),
		GoogleSearchCX:  getEnv("GOOGLE_SEARCH_CX", ""),
	}

	var err error
	cfg.JWTAccessTTL, err = strconv.Atoi(getEnv("JWT_ACCESS_TTL", "3600"))
	if err != nil {
		return nil, fmt.Errorf("JWT_ACCESS_TTL: %w", err)
	}
	cfg.JWTRefreshTTL, err = strconv.Atoi(getEnv("JWT_REFRESH_TTL", "604800"))
	if err != nil {
		return nil, fmt.Errorf("JWT_REFRESH_TTL: %w", err)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
