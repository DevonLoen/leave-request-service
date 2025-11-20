package constant

import "time"

var EnvKeys = envKeys{
	Env:                "ENV",
	ServerAddress:      "SERVER_ADDRESS",
	CorsAllowedOrigin:  "CORS_ALLOWED_ORIGIN",
	DBDriver:           "DB_DRIVER",
	DBSource:           "DB_SOURCE",
	SuperAdminEmail:    "SUPER_ADMIN_EMAIL",
	SuperAdminPassword: "SUPER_ADMIN_PASSWORD",
}

var Headers = headers{
	Origin:        "Origin",
	ContentLength: "Content-Length",
}

var MaxAge = 12 * time.Hour

type envKeys struct {
	Env                string
	ServerAddress      string
	CorsAllowedOrigin  string
	DBDriver           string
	DBSource           string
	SuperAdminEmail    string
	SuperAdminPassword string
}

type headers struct {
	Origin        string
	ContentLength string
}
