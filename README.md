# mgc-infra-costs

Estimativa de custos de infraestrutura Magalu Cloud a partir de planos Terraform.

Busca preços em tempo real da API da calculadora MGC e gera tabelas no terminal ou JSON compatível com Infracost para integração CI/CD.

## Instalação

```bash
go install github.com/prenansantana/mgc-infra-costs/cmd/mgc-infra-costs@latest
```

Ou compile a partir do código fonte:

```bash
git clone https://github.com/prenansantana/mgc-infra-costs.git
cd mgc-infra-costs
go build -o mgc-infra-costs ./cmd/mgc-infra-costs
```

## Uso

```bash
# Gerar JSON do plano Terraform
terraform plan -out=tfplan
terraform show -json tfplan > plan.json

# Saída em tabela (padrão)
mgc-infra-costs breakdown --plan plan.json

# JSON compatível com Infracost
mgc-infra-costs breakdown --plan plan.json --format json

# Especificar região
mgc-infra-costs breakdown --plan plan.json --region br-se1
```

### Exemplo de saída

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
   └─ Object Storage Standard (usage-based, R$ 0.10/GB/mês)                                  R$ 0.00
 ────────────────────────────────────────────────────────────                                  ───────────────
 TOTAL MONTHLY COST                                                                           R$ 347.20
```

## Recursos suportados

| Recurso Terraform | Descrição |
|---|---|
| `mgc_virtual_machine_instances` | Máquinas Virtuais (todos os flavors BV/DP) |
| `mgc_block_storage_volumes` | Block Storage (NVMe 1K/5K/etc) |
| `mgc_dbaas_instances` | Database as a Service (MySQL, PostgreSQL) |
| `mgc_network_public_ips` | IPs Públicos (sem custo de alocação) |
| `mgc_object_storage_buckets` | Object Storage (custo por uso) |

## Formatos de machine type

Ambos os formatos são suportados:

- Formato curto: `BV2-4-20` (classe + vCPU + RAM + disco)
- Flavor name: `i1-c2-r4-d20`

## Fonte de preços

Os preços são buscados em tempo real da API da calculadora Magalu Cloud:

```
GET https://calculadora.magalu.cloud/api/sku/v0/skus/region={region}.json
```

Não requer autenticação. Preços em BRL com impostos inclusos.

## Integração CI/CD

A saída `--format json` é compatível com o schema JSON do Infracost (`version: 0.2`), permitindo integração com ferramentas de CI/CD que consomem output do Infracost.

## Licença

MIT
