package config

type Config struct {
	Port             string
	MongoURI         string
	MongoDB          string
	JWTSecret        string
	JWTAccessTTL     int
	JWTRefreshTTL    int
	GCSBucket        string
	GCPProjectID     string
	VisionLocation   string
	YouTubeAPIKey    string
	GoogleSearchKey  string
	GoogleSearchCX   string
}

func Load() (*Config, error) {
	return nil, nil
}
