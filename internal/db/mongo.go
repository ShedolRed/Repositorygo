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

// Важно: Railway Mongo часто требует authSource=admin.
// Также: если в URI есть query (?opts), то перед ним должен быть "/" -> "/?opts".
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

	// url.URL при наличии query уже формирует "/?..." корректно
	return u.String(), nil
}

func normalizeMongoURI(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return s
	}

	// srv не трогаем (обычно корректный)
	if strings.HasPrefix(s, "mongodb+srv://") {
		return s
	}

	// 1) Если есть "?" — убедимся, что перед ним есть "/"
	qPos := strings.Index(s, "?")
	if qPos != -1 {
		// ищем первый "/" после "://"
		schemeIdx := strings.Index(s, "://")
		afterScheme := 0
		if schemeIdx != -1 {
			afterScheme = schemeIdx + 3
		}
		slashPos := strings.Index(s[afterScheme:], "/")
		if slashPos == -1 || afterScheme+slashPos > qPos {
			// нет "/" перед query -> вставляем
			s = s[:qPos] + "/" + s[qPos:]
			qPos++ // сдвиг
		}
	} else {
		// 2) Нет query — если вообще нет "/" после хоста, добавим "/" чтобы потом было "/?..."
		schemeIdx := strings.Index(s, "://")
		afterScheme := 0
		if schemeIdx != -1 {
			afterScheme = schemeIdx + 3
		}
		slashPos := strings.Index(s[afterScheme:], "/")
		if slashPos == -1 {
			s += "/"
		}
	}

	// 3) Добавим authSource=admin если его нет
	if !strings.Contains(s, "authSource=") {
		if strings.Contains(s, "?") {
			s += "&authSource=admin"
		} else {
			s += "?authSource=admin"
		}
	}

	return s
}
