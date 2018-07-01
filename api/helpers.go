package api

import (
	"net/url"
	"os"
)

// GetURL ...
func GetURL(path string) *url.URL {
	u := &url.URL{
		Scheme: "http",
		Host:   os.Getenv(EnvSonarrURL),
		Path:   path,
	}
	q := u.Query()
	q.Set("apikey", os.Getenv(EnvSonarrAPIKey))
	u.RawQuery = q.Encode()
	return u
}
