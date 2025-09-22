package env

import (
	"github.com/go-gorote/gorote"
	_ "github.com/joho/godotenv/autoload"
)

type APP struct {
	Name     string
	Version  string
	TimeZone string
	Port     int
}

var App = APP{
	Name:     gorote.MustEnvAsString("APP_NAME"),
	Version:  gorote.MustEnvAsString("APP_VERSION"),
	TimeZone: gorote.MustEnvAsString("APP_TIMEZONE", "America/Fortaleza"),
	Port:     gorote.MustEnvAsInt("APP_PORT"),
}

var Domain = gorote.MustEnvAsString("DOMAIN")
var CORS = gorote.MustEnvAsString("CORS")

var Sql01 = gorote.InitPostgres{
	Host:     gorote.MustEnvAsString("SQL_HOST", "localhost"),
	User:     gorote.MustEnvAsString("SQL_USERNAME"),
	Password: gorote.MustEnvAsString("SQL_PASSWORD"),
	Database: gorote.MustEnvAsString("SQL_DATABASE"),
	Port:     gorote.MustEnvAsInt("SQL_PORT"),
	TimeZone: gorote.MustEnvAsString("APP_TIMEZONE", "America/Fortaleza"),
	Schema:   gorote.MustEnvAsString("SQL_SCHEMA"),
}

var JwtExpireAccess = gorote.MustEnvAsTime("JWT_EXPIRE_ACCESS")
var JwtExpireRefresh = gorote.MustEnvAsTime("JWT_EXPIRE_REFRESH")

var SuperEmail = gorote.MustEnvAsString("SUPER_EMAIL")
var SuperPassword = gorote.MustEnvAsString("SUPER_PASSWORD")

var CollectorOpenTelemetry = gorote.MustEnvAsString("COLLECTOR_OPENTELEMETRY")

var PrivateKey = gorote.MustEnvPrivateKeyRSA("PRIVATE_RSA")

var PublicKey = gorote.MustEnvPublicKeyRSA("PUBLIC_RSA")
