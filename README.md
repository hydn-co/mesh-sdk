# Mesh SDK

A small Go SDK used by Hydden's mesh services and tooling. It provides common packages and helpers used across the mesh projects (authentication helpers, local store utilities, lifecycle helpers, test helpers, and more).

## Contents

- `pkg/` — Core packages (auth, secrets, tenant, lifecycle, localstore, testkit, etc.)
- `docs/` — Development and testing guidelines

## Requirements

- Go 1.25+ (ensure your `GOPATH`/modules are configured)

## Quick start

Clone the repo and run the tests or build the packages:

```shell
git clone https://github.com/hydn-co/mesh-sdk.git
cd mesh/mesh-sdk
go test ./... -v
go build ./...
```

Run the unit tests (verbose):

```shell
go test ./... -v
```

Run vet/static checks (optional):

```shell
go vet ./...
```

## Contributing

See `docs/code_guidelines.md` and `docs/test_guidelines.md` for repository conventions and testing guidance. Keep changes small and add unit tests for new behavior.
