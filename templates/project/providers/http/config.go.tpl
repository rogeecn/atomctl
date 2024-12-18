package http

import (
	"fmt"
)

const DefaultPrefix = "Http"

type Config struct {
	StaticPath  *string
	StaticRoute *string
	BaseURI     *string
	Port        uint
	Tls         *Tls
	Cors        *Cors
}

type Tls struct {
	Cert string
	Key  string
}

type Cors struct {
	Mode      string
	Whitelist []Whitelist
}

type Whitelist struct {
	AllowOrigin      string
	AllowHeaders     string
	AllowMethods     string
	ExposeHeaders    string
	AllowCredentials bool
}

func (h *Config) Address() string {
	return fmt.Sprintf(":%d", h.Port)
}
