# helmfile

Docker image for Kubernetes deploys in CI/CD pipelines.

Bundled tools:

- helmfile, helm, helm diff (from the official helmfile base image, which also ships helm-secrets, helm-s3, helm-git, kustomize, sops, age)
- kubectl (latest stable, layered on top)
- node (npm, npx included)

## Image

Published to GHCR as `ghcr.io/alexeyco/helmfile`:

- Tags: `latest` and the bundled helmfile version (e.g. `1.7.4`)
- Multi-arch: `linux/amd64` and `linux/arm64`
- Alpine-based via the official helmfile base image (shell available, so
  scripts work)
- Rebuilt daily at 03:00 UTC via CI (also manually triggerable); untagged
  versions are cleaned up automatically

## Development

Requirements: Go 1.27+, golangci-lint v2+.

| Command         | Description                                   |
| --------------- | --------------------------------------------- |
| `make`          | Show help                                     |
| `make mod`      | Tidy and download dependencies                |
| `make fmt`      | Format code (golangci-lint)                   |
| `make lint`     | Run linters (golangci-lint)                   |
| `make test`     | Run tests                                     |
| `make build`    | Build binary into `./build`                   |
| `make run`      | Run the app                                   |
| `make generate` | Generate Dockerfile with latest tool versions |
| `make image`    | Generate Dockerfile and build the local image |
| `make mock`     | Generate mocks (uber/mock)                    |

## License

[BSD-2-Clause](LICENSE)
