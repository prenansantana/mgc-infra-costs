package pricing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultBaseURL = "https://calculadora.magalu.cloud/api/sku/v0/skus"

type RAM struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type Disk struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type SKU struct {
	PricePerUnit    string      `json:"price_per_unit"`
	Description     string      `json:"description"`
	PriceUnit       string      `json:"price_unit"`
	SKUExternalName string      `json:"sku_external_name"`
	FlavorName      string      `json:"flavor_name"`
	HourlyPrice     string      `json:"hourly_price"`
	MonthlyPrice    string      `json:"monthly_price"`
	Region          string      `json:"region"`
	SKUClass        string      `json:"sku_class"`
	VCPUs           int         `json:"vcpu"`
	RAM             *RAM        `json:"ram"`
	Disk            *Disk       `json:"disk"`
	Architecture    string      `json:"architecture"`
	GPU             interface{} `json:"gpu"`
	Availability    string      `json:"availability"`
	System          string      `json:"system"`
	Class           string      `json:"class"`
	RelatedSKUs     []string    `json:"related_skus_external_names"`
}

type ProductCategory struct {
	Category string `json:"category"`
	Product  string `json:"product"`
	SKUs     []SKU  `json:"skus"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) FetchPricing(region string) ([]ProductCategory, error) {
	url := fmt.Sprintf("%s/region=%s.json", c.baseURL, region)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching pricing data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pricing API returned status %d", resp.StatusCode)
	}

	var categories []ProductCategory
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("decoding pricing data: %w", err)
	}

	return categories, nil
}
