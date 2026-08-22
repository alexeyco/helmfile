package generate_test

import (
	"strings"
	"testing"

	"github.com/alexeyco/helmfile/internal"
	"github.com/alexeyco/helmfile/internal/generate"
	"github.com/spf13/afero"
)

func testVersions() internal.Versions {
	return internal.Versions{
		Alpine:   "3.24",
		Kubectl:  "1.36.4",
		Helm:     "3.21.4",
		Helmfile: "1.7.4",
		HelmDiff: "3.15.11",
	}
}

func TestRender(t *testing.T) {
	g := generate.New(afero.NewMemMapFs())

	got, err := g.Render(testVersions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{
		"FROM alpine:3.24 AS builder",
		"KUBECTL_VERSION=1.36.4",
		"HELM_VERSION=3.21.4",
		"HELMFILE_VERSION=1.7.4",
		"HELMDIFF_VERSION=3.15.11",
		"helm-v${HELM_VERSION}-linux-${TARGETARCH}.tar.gz",
		"helmfile_${HELMFILE_VERSION}_linux_${TARGETARCH}.tar.gz",
		"helm-diff-linux-${TARGETARCH}.tgz",
		"distroless/static-debian12",
		"HELM_DATA_HOME",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestGenerate(t *testing.T) {
	fs := afero.NewMemMapFs()
	g := generate.New(fs)

	if err := g.Generate(testVersions(), "Dockerfile"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := afero.ReadFile(fs, "Dockerfile")
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	if !strings.Contains(string(content), "FROM alpine:3.24 AS builder") {
		t.Errorf("generated file missing alpine base")
	}
}
