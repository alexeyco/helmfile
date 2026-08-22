package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/alexeyco/helmfile/internal"
)

const (
	kubectlURL  = "https://dl.k8s.io/release/stable.txt"
	helmURL     = "https://api.github.com/repos/helm/helm/releases/latest"
	helmfileURL = "https://api.github.com/repos/helmfile/helmfile/releases/latest"
	helmDiffURL = "https://api.github.com/repos/databus23/helm-diff/releases/latest"
	alpineURL   = "https://dl-cdn.alpinelinux.org/alpine/"

	githubAccept = "application/vnd.github+json"
	userAgent    = "github.com/alexeyco/helmfile (image generator)"
)

var githubHeader = map[string]string{
	"Accept":     githubAccept,
	"User-Agent": userAgent,
}

var alpineReleaseRe = regexp.MustCompile(`v(3\.\d+)/`)

var _ HTTPClient = (*http.Client)(nil)

type Fetcher struct {
	client HTTPClient
}

func New(client HTTPClient) *Fetcher {
	return &Fetcher{client: client}
}

type source struct {
	name   string
	url    string
	header map[string]string
	parse  func(body []byte) (string, error)
}

func (f *Fetcher) Fetch(ctx context.Context) (internal.Versions, error) {
	sources := []source{
		{name: "kubectl", url: kubectlURL, parse: parsePlainVersion},
		{name: "helm", url: helmURL, header: githubHeader, parse: parseTagName},
		{name: "helmfile", url: helmfileURL, header: githubHeader, parse: parseTagName},
		{name: "helm-diff", url: helmDiffURL, header: githubHeader, parse: parseTagName},
		{name: "alpine", url: alpineURL, parse: parseAlpineIndex},
	}

	results := make([]string, len(sources))
	errs := make([]error, len(sources))

	var wg sync.WaitGroup
	for i, s := range sources {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := f.fetchSource(ctx, s)
			if err != nil {
				err = fmt.Errorf("%s: %w", s.name, err)
			}
			results[i], errs[i] = v, err
		}()
	}
	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return internal.Versions{}, fmt.Errorf("fetch versions: %w", err)
	}

	return internal.Versions{
		Kubectl:  results[0],
		Helm:     results[1],
		Helmfile: results[2],
		HelmDiff: results[3],
		Alpine:   results[4],
	}, nil
}

func (f *Fetcher) fetchSource(ctx context.Context, s source) (string, error) {
	body, err := f.get(ctx, s)
	if err != nil {
		return "", err
	}
	return s.parse(body)
}

func (f *Fetcher) get(ctx context.Context, s source) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	for k, v := range s.header {
		req.Header.Set(k, v)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return body, nil
}

func parsePlainVersion(body []byte) (string, error) {
	v, err := validateVersion(string(body))
	if err != nil {
		return "", fmt.Errorf("parse plain version: %w", err)
	}
	return v, nil
}

func parseTagName(body []byte) (string, error) {
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("decode JSON: %w", err)
	}
	v, err := validateVersion(payload.TagName)
	if err != nil {
		return "", fmt.Errorf("parse tag_name: %w", err)
	}
	return v, nil
}

func parseAlpineIndex(body []byte) (string, error) {
	latest := ""
	for _, m := range alpineReleaseRe.FindAllStringSubmatch(string(body), -1) {
		if compareMinor(m[1], latest) > 0 {
			latest = m[1]
		}
	}
	if latest == "" {
		return "", errors.New("no alpine release found")
	}
	return latest, nil
}

func compareMinor(a, b string) int {
	if b == "" {
		return 1
	}
	pa := strings.Split(a, ".")
	pb := strings.Split(b, ".")
	for i := range 2 {
		var na, nb int
		if i < len(pa) {
			na, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			nb, _ = strconv.Atoi(pb[i])
		}
		if na != nb {
			return na - nb
		}
	}
	return 0
}

func validateVersion(raw string) (string, error) {
	v := strings.TrimSpace(raw)
	v = strings.TrimPrefix(strings.TrimPrefix(v, "v"), "V")
	if v == "" {
		return "", errors.New("empty version")
	}
	if v == "latest" {
		return "", errors.New(`invalid version "latest"`)
	}
	return v, nil
}
