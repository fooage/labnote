package config

import (
	"time"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/spf13/viper"
)

// LoadConfiguration read the server-related setting information from the configuration file.
func LoadConfiguration() (data.Database, cache.Cache, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, err
	}
	// Read in database and cache setting information.
	databaseType := viper.GetString("database.type")
	databaseCommand := viper.GetString("database.command")
	databaseName := viper.GetString("database.name")
	cacheAddress := viper.GetString("cache.address")
	cachePassword := viper.GetString("cache.password")
	cachePoolSize := viper.GetInt("cache.pool_size")
	// Settings, such as cookie expiration time, scope, url prefix, token settings.
	handler.CookieExpireDuration = viper.GetInt("server.handler.cookie_duration")
	handler.CookieAccessScope = viper.GetString("server.handler.access_scope")
	handler.FileStorageDirectory = viper.GetString("server.handler.storage_directory")
	handler.DownloadUrlBase = viper.GetString("server.handler.url_base")
	handler.TokenExpireDuration = time.Second * time.Duration(viper.GetInt64("server.token.token_duration"))
	handler.EncryptionKey = viper.GetString("server.token.encryption_key")
	handler.TokenIssuer = viper.GetString("server.token.token_issuer")
	// Generate database and cache database objects and return.
	var db data.Database
	var ch cache.Cache
	switch databaseType {
	case "mongo":
		db = data.NewMongoDB(databaseCommand, databaseName)
	case "mysql":
		db = data.NewMySQL(databaseCommand, databaseName)
	default:
		db = nil
	}
	ch = cache.NewRedis(cacheAddress, cachePassword, cachePoolSize)
	return db, ch, nil
}
