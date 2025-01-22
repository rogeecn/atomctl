package req

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"{{.ModuleName}}/providers/req/cookiejar"

	"github.com/imroc/req/v3"
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"
)

type Client struct {
	client *req.Client
	jar    *cookiejar.Jar
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func() (*Client, error) {
		c := &Client{}

		client := req.C()
		if config.DevMode {
			client.DevMode()
		}

		if config.CookieJarFile != "" {
			dir := filepath.Dir(config.CookieJarFile)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0o755)
				if err != nil {
					return nil, err
				}
			}
			jar, err := cookiejar.New(&cookiejar.Options{
				Filename: config.CookieJarFile,
			})
			if err != nil {
				return nil, err
			}

			c.jar = jar
			client.SetCookieJar(jar)
		}

		if config.RootCa != nil {
			client.SetRootCertsFromFile(config.RootCa...)
		}

		if config.InsecureSkipVerify {
			client.EnableInsecureSkipVerify()
		}

		if config.UserAgent != "" {
			client.SetUserAgent(config.UserAgent)
		}
		if config.Timeout > 0 {
			client.SetTimeout(time.Duration(config.Timeout) * time.Second)
		}

		if config.CommonHeaders != nil {
			client.SetCommonHeaders(config.CommonHeaders)
		}

		if config.AuthBasic.Username != "" && config.AuthBasic.Password != "" {
			client.SetCommonBasicAuth(config.AuthBasic.Username, config.AuthBasic.Password)
		}

		if config.AuthBearerToken != "" {
			client.SetCommonBearerAuthToken(config.AuthBearerToken)
		}

		if config.ProxyURL != "" {
			client.SetProxyURL(config.ProxyURL)
		}

		if config.RedirectPolicy != nil {
			client.SetRedirectPolicy(parsePolicies(config.RedirectPolicy)...)
		}

		c.client = client
		return c, nil
	}, o.DiOptions()...)
}

func parsePolicies(policies []string) []req.RedirectPolicy {
	ps := []req.RedirectPolicy{}
	for _, policy := range policies {
		policyItems := strings.Split(policy, ":")
		if len(policyItems) != 2 {
			continue
		}

		switch policyItems[0] {
		case "Max":
			max, err := strconv.Atoi(policyItems[1])
			if err != nil {
				continue
			}
			ps = append(ps, req.MaxRedirectPolicy(max))
		case "No":
			ps = append(ps, req.NoRedirectPolicy())
		case "SameDomain":
			ps = append(ps, req.SameDomainRedirectPolicy())
		case "SameHost":
			ps = append(ps, req.SameHostRedirectPolicy())
		case "AllowedHost":
			ps = append(ps, req.AllowedHostRedirectPolicy(strings.Split(policyItems[1], ",")...))
		case "AllowedDomain":
			ps = append(ps, req.AllowedDomainRedirectPolicy(strings.Split(policyItems[1], ",")...))
		}
	}

	return ps
}

func (c *Client) R() *req.Request {
	return c.client.R()
}

func (c *Client) SaveCookJar() error {
	return c.jar.Save()
}

func (c *Client) GetCookie(key string) (string, bool) {
	kv := c.AllCookiesKV()
	v, ok := kv[key]
	return v, ok
}

func (c *Client) AllCookies() []*http.Cookie {
	return c.jar.AllCookies()
}

func (c *Client) AllCookiesKV() map[string]string {
	return c.jar.KVData()
}

func (c *Client) SetCookie(u *url.URL, cookies []*http.Cookie) {
	c.jar.SetCookies(u, cookies)
}
