package fetcher_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexeyco/helmfile/internal"
	"github.com/alexeyco/helmfile/internal/fetcher"
	"go.uber.org/mock/gomock"
)

const alpineIndexBody = `<html>
<head><title>Index of /alpine/</title></head>
<body>
<h1>Index of /alpine/</h1><hr><pre><a href="../">../</a>
<a href="edge/">edge/</a>
<a href="latest-stable/">latest-stable/</a>
<a href="v3.20/">v3.20/</a>
<a href="v3.24/">v3.24/</a>
<a href="v3.9/">v3.9/</a>
</pre><hr></body>
</html>`

func TestFetcherFetch(t *testing.T) {
	responses := map[string]string{
		"https://dl.k8s.io/release/stable.txt":                             "v1.36.4\n",
		"https://api.github.com/repos/helm/helm/releases/latest":           `{"tag_name": "v3.21.4"}`,
		"https://api.github.com/repos/helmfile/helmfile/releases/latest":   `{"tag_name": "v1.7.4"}`,
		"https://api.github.com/repos/databus23/helm-diff/releases/latest": `{"tag_name": "v3.15.11"}`,
		"https://dl-cdn.alpinelinux.org/alpine/":                           alpineIndexBody,
	}

	ctrl := gomock.NewController(t)
	client := NewMockHTTPClient(ctrl)
	client.EXPECT().Do(gomock.Any()).Times(5).DoAndReturn(
		func(req *http.Request) (*http.Response, error) {
			body, ok := responses[req.URL.String()]
			if !ok {
				t.Errorf("unexpected url %s", req.URL)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		},
	)

	got, err := fetcher.New(client).Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := internal.Versions{
		Alpine:   "3.24",
		Kubectl:  "1.36.4",
		Helm:     "3.21.4",
		Helmfile: "1.7.4",
		HelmDiff: "3.15.11",
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestFetcherFetchHTTPError(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockHTTPClient(ctrl)
	client.EXPECT().Do(gomock.Any()).Times(5).Return(nil, errors.New("boom"))

	if _, err := fetcher.New(client).Fetch(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}
