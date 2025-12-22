package config

import (
    "os"
    "strings"
)

type Config struct {
    Port string

    MongoDBName string
    MongoURL    string

    MongoHost     string
    MongoPort     string
    MongoUser     string
    MongoPassword string

    MongoRootUser     string
    MongoRootPassword string
}

func Load() Config {
    return Config{
        Port:        env("PORT", "8080"),
        MongoDBName: env("MONGO_DB_NAME", "appdb"),

        MongoURL: strings.TrimSpace(os.Getenv("MONGO_URL")),

        MongoHost:     strings.TrimSpace(os.Getenv("MONGOHOST")),
        MongoPort:     strings.TrimSpace(os.Getenv("MONGOPORT")),
        MongoUser:     strings.TrimSpace(os.Getenv("MONGOUSER")),
        MongoPassword: strings.TrimSpace(os.Getenv("MONGOPASSWORD")),

        MongoRootUser:     strings.TrimSpace(os.Getenv("MONGO_INITDB_ROOT_USERNAME")),
        MongoRootPassword: strings.TrimSpace(os.Getenv("MONGO_INITDB_ROOT_PASSWORD")),
    }
}

func env(key, def string) string {
    v := strings.TrimSpace(os.Getenv(key))
    if v == "" {
        return def
    }
    return v
}
