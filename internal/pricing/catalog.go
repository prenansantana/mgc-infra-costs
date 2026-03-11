package pricing

import (
	"fmt"
	"strconv"
)

type Catalog struct {
	byProduct map[string]map[string]SKU // product -> flavor_name -> SKU
}

func NewCatalog(categories []ProductCategory) *Catalog {
	c := &Catalog{
		byProduct: make(map[string]map[string]SKU),
	}
	for _, cat := range categories {
		if _, ok := c.byProduct[cat.Product]; !ok {
			c.byProduct[cat.Product] = make(map[string]SKU)
		}
		for _, sku := range cat.SKUs {
			c.byProduct[cat.Product][sku.FlavorName] = sku
		}
	}
	return c
}

func (c *Catalog) LookupByFlavor(product, flavorName string) (*SKU, error) {
	products, ok := c.byProduct[product]
	if !ok {
		return nil, fmt.Errorf("product %q not found in catalog", product)
	}
	sku, ok := products[flavorName]
	if !ok {
		return nil, fmt.Errorf("flavor %q not found in product %q", flavorName, product)
	}
	return &sku, nil
}

func (c *Catalog) ListFlavors(product string) []SKU {
	products, ok := c.byProduct[product]
	if !ok {
		return nil
	}
	skus := make([]SKU, 0, len(products))
	for _, sku := range products {
		skus = append(skus, sku)
	}
	return skus
}

// LookupBySpecs finds a VM SKU matching the given class prefix, vcpu, ram and disk.
// This handles the Terraform machine_type format like "BV2-4-20" which maps to
// sku_class="Balanced Value (BV)", vcpu=2, ram=4GB, disk=20GB.
func (c *Catalog) LookupBySpecs(product, classPrefix string, vcpu, ramGB, diskGB int) (*SKU, error) {
	products, ok := c.byProduct[product]
	if !ok {
		return nil, fmt.Errorf("product %q not found in catalog", product)
	}
	for _, sku := range products {
		if sku.VCPUs != vcpu || sku.RAM == nil || sku.Disk == nil {
			continue
		}
		if sku.RAM.Value != ramGB || sku.Disk.Value != diskGB {
			continue
		}
		// Match class prefix (e.g., "BV" matches "Balanced Value (BV)")
		if classPrefix != "" {
			classTag := fmt.Sprintf("(%s)", classPrefix)
			if !containsStr(sku.SKUClass, classTag) {
				continue
			}
		}
		return &sku, nil
	}
	return nil, fmt.Errorf("no SKU found in %q matching class=%s vcpu=%d ram=%dGB disk=%dGB", product, classPrefix, vcpu, ramGB, diskGB)
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && searchStr(s, substr)))
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (c *Catalog) Products() []string {
	products := make([]string, 0, len(c.byProduct))
	for p := range c.byProduct {
		products = append(products, p)
	}
	return products
}

func MonthlyPriceFloat(sku *SKU) float64 {
	price, err := strconv.ParseFloat(sku.MonthlyPrice, 64)
	if err != nil {
		return 0
	}
	return price
}

func HourlyPriceFloat(sku *SKU) float64 {
	price, err := strconv.ParseFloat(sku.HourlyPrice, 64)
	if err != nil {
		return 0
	}
	return price
}
