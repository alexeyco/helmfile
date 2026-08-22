package fetcher

import "net/http"

//go:generate go tool mockgen -source=deps.go -destination=deps_test.go -package=fetcher_test

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
