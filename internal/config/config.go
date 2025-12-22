package config

import "os"

type Config struct {
	Port       string
	MongoURL   string
	MongoDBName string
}

func Load() Config {
	port := getenv("PORT", "8080")

	// Railway Mongo service provides MONGO_URL
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		// allow alternative naming for local/other hosts
		mongoURL = os.Getenv("MONGODB_URL")
	}

	dbName := getenv("MONGO_DB_NAME", "appdb")

	return Config{
		Port:        port,
		MongoURL:    mongoURL,
		MongoDBName: dbName,
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
