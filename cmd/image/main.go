package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexeyco/helmfile/internal/fetcher"
	"github.com/alexeyco/helmfile/internal/generate"
	"github.com/spf13/afero"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	versions, err := fetcher.New(client).Fetch(ctx)
	if err != nil {
		return fmt.Errorf("fetch versions: %w", err)
	}

	if err := generate.New(afero.NewOsFs()).Generate(versions, "Dockerfile"); err != nil {
		return fmt.Errorf("generate Dockerfile: %w", err)
	}

	fmt.Printf("kubectl   %s\n", versions.Kubectl)
	fmt.Printf("helm      %s\n", versions.Helm)
	fmt.Printf("helmfile  %s\n", versions.Helmfile)
	fmt.Printf("helm-diff %s\n", versions.HelmDiff)
	fmt.Printf("alpine    %s\n", versions.Alpine)
	fmt.Println("Dockerfile written")
	return nil
}
