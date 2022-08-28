package config

import (
	"time"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/fooage/labnote/proxy"
	"github.com/spf13/viper"
)

// ServerConfiguration read the server-related setting information from the configuration file.
func ServerConfiguration() (data.Database, cache.Cache, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, err
	}

	// Read in database and cache setting information. Settings, such as cookie
	// expiration time, scope, url prefix, token settings.
	handler.HostAddress = viper.GetString("server.handler.host_address")
	handler.ListenPort = viper.GetString("server.handler.listen_port")
	handler.CookieExpireDuration = viper.GetInt("server.handler.cookie_duration")
	handler.CookieAccessScope = viper.GetString("server.handler.access_scope")
	handler.FileStorageDirectory = viper.GetString("server.handler.storage_directory")
	handler.DownloadUrlBase = viper.GetString("server.handler.url_base")
	handler.TokenExpireDuration = time.Second * time.Duration(viper.GetInt64("server.token.token_duration"))
	handler.EncryptionKey = viper.GetString("server.token.encryption_key")
	handler.TokenIssuer = viper.GetString("server.token.token_issuer")

	// The Settings of proxy server's connect information.
	proxy.ProxyAddress = viper.GetString("server.proxy.proxy_address")
	proxy.UseProxy = viper.GetBool("server.proxy.use_proxy")
	proxy.Timeout = time.Second * time.Duration(viper.GetInt64("server.proxy.timeout"))

	// Generate database and cache database objects and return.
	var db data.Database
	var ch cache.Cache
	switch viper.GetString("database.type") {
	case "mongo":
		db = data.NewMongoDB(viper.GetString("database.command"), viper.GetString("database.name"))
	case "mysql":
		db = data.NewMySQL(viper.GetString("database.command"), viper.GetString("database.name"))
	default:
		db = nil
	}
	ch = cache.NewRedis(viper.GetString("cache.address"), viper.GetString("cache.password"), viper.GetInt("cache.pool_size"))
	return db, ch, nil
}

// Load the configuration information of the reverse proxy.
func ProxyConfiguration() (data.Database, cache.Cache, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, err
	}

	// Read in database and cache setting information. Reverse proxy server configuration.
	proxy.HostAddress = viper.GetString("proxy.universal.host_address")
	proxy.ListenPort = viper.GetString("proxy.universal.listen_port")
	proxy.Timeout = time.Second * time.Duration(viper.GetInt64("proxy.registry.timeout"))
	proxy.HttpVersion = viper.GetString("proxy.universal.http_version")

	// Generate database and cache database objects and return.
	var db data.Database
	var ch cache.Cache
	switch viper.GetString("database.type") {
	case "mongo":
		db = data.NewMongoDB(viper.GetString("database.command"), viper.GetString("database.name"))
	case "mysql":
		db = data.NewMySQL(viper.GetString("database.command"), viper.GetString("database.name"))
	default:
		db = nil
	}
	ch = cache.NewRedis(viper.GetString("cache.address"), viper.GetString("cache.password"), viper.GetInt("cache.pool_size"))
	return db, ch, nil
}
