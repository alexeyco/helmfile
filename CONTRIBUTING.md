# Contributing

## Branching

- Never push to `master`.
- Branch from `master`: `git checkout -b feat/<topic>`.
- Commit messages: no conventional prefixes (`feat:`, `fix:`, …), past tense —
  e.g. `Added node to the image`, `Fixed the CI build`.

## Development

Requirements and the full list of Makefile targets: see
[Development](README.md#development) in the README.

Work through the Makefile — never invoke `go` or `golangci-lint` directly.
After changing code, run in this order:

```sh
make mod fmt lint test
```

Notes:

- Regenerate mocks after changing fetcher interfaces: `make mock`.
- `make generate` writes the root `Dockerfile` (gitignored); it is produced
  by CI and `make image`, never committed.
- Format Markdown with `npx prettier --write <file>.md`.

## Publishing

Nothing to do manually. The `docker` workflow
([`.github/workflows/docker.yml`](.github/workflows/docker.yml)) builds the
multi-arch image and pushes it to GHCR daily at 03:00 UTC, or on demand via
`workflow_dispatch`. Untagged image versions are cleaned up automatically.
There are no releases, tags or changelog — image tags
(`latest` + the bundled helmfile version) are derived at build time.
