package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/pelletier/go-toml/v2"
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

// tomlFile mirrors config.toml structure.
type tomlFile struct {
	Server struct {
		Port string `toml:"port"`
	} `toml:"server"`
	Mongo struct {
		URI string `toml:"uri"`
		DB  string `toml:"db"`
	} `toml:"mongo"`
	JWT struct {
		Secret     string `toml:"secret"`
		AccessTTL  int    `toml:"access_ttl"`
		RefreshTTL int    `toml:"refresh_ttl"`
	} `toml:"jwt"`
	GCP struct {
		ProjectID      string `toml:"project_id"`
		GCSBucket      string `toml:"gcs_bucket"`
		VisionLocation string `toml:"vision_location"`
	} `toml:"gcp"`
	External struct {
		YouTubeAPIKey   string `toml:"youtube_api_key"`
		GoogleSearchKey string `toml:"google_search_key"`
		GoogleSearchCX  string `toml:"google_search_cx"`
	} `toml:"external"`
}

// Load reads config.toml from the project root, then overrides with env vars
// (env vars take precedence so CI/CD secrets work without modifying the file).
func Load() (*Config, error) {
	tf := defaultToml()
	if data, err := os.ReadFile(findConfigFile()); err == nil {
		_ = toml.Unmarshal(data, &tf) // ignore parse errors — fall back to defaults
	}

	cfg := &Config{
		Port:            or(os.Getenv("PORT"), tf.Server.Port),
		MongoURI:        or(os.Getenv("MONGODB_URI"), tf.Mongo.URI),
		MongoDB:         or(os.Getenv("MONGODB_DB"), tf.Mongo.DB),
		JWTSecret:       or(os.Getenv("JWT_SECRET"), tf.JWT.Secret),
		JWTAccessTTL:    tf.JWT.AccessTTL,
		JWTRefreshTTL:   tf.JWT.RefreshTTL,
		GCSBucket:       or(os.Getenv("GCS_BUCKET"), tf.GCP.GCSBucket),
		GCPProjectID:    or(os.Getenv("GCP_PROJECT_ID"), tf.GCP.ProjectID),
		VisionLocation:  or(os.Getenv("VISION_LOCATION"), tf.GCP.VisionLocation),
		YouTubeAPIKey:   or(os.Getenv("YOUTUBE_API_KEY"), tf.External.YouTubeAPIKey),
		GoogleSearchKey: or(os.Getenv("GOOGLE_SEARCH_KEY"), tf.External.GoogleSearchKey),
		GoogleSearchCX:  or(os.Getenv("GOOGLE_SEARCH_CX"), tf.External.GoogleSearchCX),
	}

	// Env var overrides for TTLs
	if v := os.Getenv("JWT_ACCESS_TTL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.JWTAccessTTL = n
		}
	}
	if v := os.Getenv("JWT_REFRESH_TTL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.JWTRefreshTTL = n
		}
	}

	return cfg, nil
}

func defaultToml() tomlFile {
	var tf tomlFile
	tf.Server.Port = "8080"
	tf.Mongo.URI = "mongodb://localhost:27017"
	tf.Mongo.DB = "OKKY"
	tf.JWT.AccessTTL = 3600
	tf.JWT.RefreshTTL = 604800
	tf.GCP.VisionLocation = "us-central1"
	return tf
}

// findConfigFile looks for config.toml relative to the binary or source root.
func findConfigFile() string {
	// 1. Beside the running binary
	if exe, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(exe), "config.toml")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// 2. Walk up from source file (convenient during `go run`)
	_, filename, _, _ := runtime.Caller(0)
	// filename = .../backend/pkg/config/config.go → go up 3 dirs to backend/
	root := filepath.Join(filepath.Dir(filename), "..", "..")
	p := filepath.Join(root, "config.toml")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	// 3. Current working directory
	return "config.toml"
}

// or returns a if non-empty, otherwise b.
func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
