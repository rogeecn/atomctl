package req

import (
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

const DefaultPrefix = "HttpClient"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	DevMode            bool
	CookieJarFile      string
	RootCa             []string
	UserAgent          string
	InsecureSkipVerify bool
	CommonHeaders      map[string]string
	Timeout            uint
	AuthBasic          struct {
		Username string
		Password string
	}
	AuthBearerToken string
	ProxyURL        string
	RedirectPolicy  []string // "Max:10;No;SameDomain;SameHost;AllowedHost:x,x,x,x,x,AllowedDomain:x,x,x,x,x"
}
