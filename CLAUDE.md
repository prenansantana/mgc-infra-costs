# CLAUDE.md

## Project

mgc-infra-costs: Go CLI that estimates Magalu Cloud infrastructure costs from Terraform plans.

## Build & Run

```bash
go build ./...                                    # build all packages
go build -o mgc-infra-costs ./cmd/mgc-infra-costs # build binary
go vet ./...                                      # lint
gofmt -w .                                        # format

# test against sample plan
go run ./cmd/mgc-infra-costs breakdown --plan testdata/sample-plan.json
go run ./cmd/mgc-infra-costs breakdown --plan testdata/sample-plan.json --format json
```

## Architecture

- `cmd/mgc-infra-costs/main.go` — CLI entrypoint (Cobra)
- `internal/pricing/client.go` — HTTP client for MGC pricing API
- `internal/pricing/catalog.go` — SKU lookup (by flavor_name or by specs)
- `internal/parser/tfplan.go` — Terraform plan JSON parser
- `internal/resources/` — Cost handlers per resource type (VM, block storage, DBaaS, network, object storage)
- `internal/output/` — Table and Infracost JSON formatters

## Key conventions

- Resource handlers implement `ResourceHandler` interface and register in `registry.go`
- Machine types support two formats: short (`BV2-4-20`) and flavor (`i1-c2-r4-d20`)
- Pricing comes from `calculadora.magalu.cloud/api/sku/v0/skus/region={region}.json`
- Output JSON follows Infracost schema v0.2
- Currency is BRL (from the API)

## Release

- GoReleaser config: `.goreleaser.yml`
- GitHub Action: `.github/workflows/release.yml` — triggers on `v*` tags
- Publishes binaries (linux/darwin × amd64/arm64) + Homebrew formula to `prenansantana/homebrew-tap`
- Requires `GORELEASER_TOKEN` secret in GitHub repo settings

To release: `git tag vX.Y.Z && git push origin vX.Y.Z`

## Adding a new resource type

1. Create `internal/resources/{name}.go` implementing `ResourceHandler`
2. Add to `supportedTypes` in `internal/parser/tfplan.go`
3. Register in `internal/resources/registry.go` init()

## Documentation rule

When making changes to features, supported resources, CLI flags, or architecture, update both:
- `README.md` (documentação para humanos, em português)
- `CLAUDE.md` (documentação para IA)
