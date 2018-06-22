package main

import (
	"fmt"
	"net/url"
	"os"
)

// GetURL ...
func GetURL(path string) *url.URL {
	fmt.Println(os.Getenv(EnvSonarrURL))
	u := &url.URL{
		Scheme: "http",
		Host:   os.Getenv(EnvSonarrURL),
		Path:   path,
	}
	fmt.Println(u.String())
	q := u.Query()
	q.Set("apikey", os.Getenv(EnvSonarrAPIKey))
	u.RawQuery = q.Encode()
	return u
}
