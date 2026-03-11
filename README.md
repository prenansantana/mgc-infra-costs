# mgc-infra-costs

Estimate Magalu Cloud infrastructure costs from Terraform plans.

Fetches live pricing from the MGC calculator API and produces terminal tables or Infracost-compatible JSON for CI/CD integration.

## Install

```bash
go install github.com/prenansantana/mgc-infra-costs/cmd/mgc-infra-costs@latest
```

Or build from source:

```bash
git clone https://github.com/prenansantana/mgc-infra-costs.git
cd mgc-infra-costs
go build -o mgc-infra-costs ./cmd/mgc-infra-costs
```

## Usage

```bash
# Generate Terraform plan JSON
terraform plan -out=tfplan
terraform show -json tfplan > plan.json

# Table output (default)
mgc-infra-costs breakdown --plan plan.json

# Infracost-compatible JSON
mgc-infra-costs breakdown --plan plan.json --format json

# Specify region
mgc-infra-costs breakdown --plan plan.json --region br-se1
```

### Example output

```
 Name                                                                                          Monthly Cost
 ────────────────────────────────────────────────────────────                                  ───────────────
 module.vm_n8n.mgc_virtual_machine_instances.this
   └─ Instance (BV2-4-20, Balanced Value (BV) - 2 vCPU, 4GB RAM, 20GB Disk)                  R$ 92.99
 module.vm_chatwoot.mgc_virtual_machine_instances.this
   └─ Instance (BV4-8-20, Balanced Value (BV) - 4 vCPU, 8GB RAM, 20GB Disk)                  R$ 159.99
 module.database.mgc_dbaas_instances.this
   └─ DBaaS PostgreSQL single_instance (1 vCPU; 4GB RAM; 10GB Disk)                          R$ 94.22
 module.storage.mgc_object_storage_buckets.this
   └─ Object Storage Standard (usage-based, R$ 0.10/GB/month)                                R$ 0.00
 ────────────────────────────────────────────────────────────                                  ───────────────
 TOTAL MONTHLY COST                                                                           R$ 347.20
```

## Supported resources

| Terraform Resource | Description |
|---|---|
| `mgc_virtual_machine_instances` | Virtual Machines (all BV/DP flavors) |
| `mgc_block_storage_volumes` | Block Storage (NVMe 1K/5K/etc) |
| `mgc_dbaas_instances` | Database as a Service (MySQL, PostgreSQL) |
| `mgc_network_public_ips` | Public IPs (no allocation charge) |
| `mgc_object_storage_buckets` | Object Storage (usage-based) |

## Machine type formats

Both formats are supported:

- Short format: `BV2-4-20` (class + vCPU + RAM + disk)
- Flavor name: `i1-c2-r4-d20`

## Pricing source

Prices are fetched live from the Magalu Cloud calculator API:

```
GET https://calculadora.magalu.cloud/api/sku/v0/skus/region={region}.json
```

No authentication required. Prices are in BRL with taxes included.

## CI/CD Integration

The `--format json` output is compatible with Infracost's JSON schema (`version: 0.2`), enabling integration with existing CI/CD tools that consume Infracost output.

## License

MIT
