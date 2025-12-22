package db

import (
    "context"
    "errors"
    "fmt"
    "net/url"
    "strings"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"

    "repositorygo/internal/config"
)

// Важно: Railway Mongo чаще всего требует authSource=admin.
// Мы добавляем его автоматически, если в URI его нет.
func ConnectMongo(ctx context.Context, cfg config.Config) (*mongo.Client, *mongo.Database, error) {
    uri, err := buildMongoURI(cfg)
    if err != nil {
        return nil, nil, err
    }

    cctx, cancel := context.WithTimeout(ctx, 12*time.Second)
    defer cancel()

    client, err := mongo.Connect(cctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, nil, fmt.Errorf("mongo connect failed: %w", err)
    }

    pctx, pcancel := context.WithTimeout(context.Background(), 6*time.Second)
    defer pcancel()
    if err := client.Ping(pctx, readpref.Primary()); err != nil {
        _ = client.Disconnect(context.Background())
        return nil, nil, fmt.Errorf("mongo ping failed: %w", err)
    }

    return client, client.Database(cfg.MongoDBName), nil
}

func buildMongoURI(cfg config.Config) (string, error) {
    if cfg.MongoURL != "" {
        return normalizeMongoURI(cfg.MongoURL), nil
    }

    host := cfg.MongoHost
    port := cfg.MongoPort
    if host == "" || port == "" {
        return "", errors.New("missing MONGO_URL (or MONGOHOST/MONGOPORT)")
    }

    user := cfg.MongoUser
    pass := cfg.MongoPassword

    if user == "" || pass == "" {
        if cfg.MongoRootUser != "" && cfg.MongoRootPassword != "" {
            user = cfg.MongoRootUser
            pass = cfg.MongoRootPassword
        }
    }
    if user == "" || pass == "" {
        return "", errors.New("missing mongo credentials: set MONGO_URL or (MONGOUSER/MONGOPASSWORD) or (MONGO_INITDB_ROOT_USERNAME/MONGO_INITDB_ROOT_PASSWORD)")
    }

    u := url.URL{
        Scheme: "mongodb",
        Host:   fmt.Sprintf("%s:%s", strings.TrimSpace(host), strings.TrimSpace(port)),
        User:   url.UserPassword(user, pass),
    }

    q := url.Values{}
    q.Set("authSource", "admin")
    q.Set("retryWrites", "false")
    u.RawQuery = q.Encode()

    return u.String(), nil
}

func normalizeMongoURI(raw string) string {
    s := strings.TrimSpace(raw)

    if strings.HasPrefix(s, "mongodb+srv://") {
        return s
    }

    if strings.Contains(s, "?") {
        if strings.Contains(s, "authSource=") {
            return s
        }
        return s + "&authSource=admin"
    }

    return s + "?authSource=admin"
}
