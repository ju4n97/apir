package manifest_test

import (
	"encoding/json"
	"testing"

	"github.com/ju4n97/esquema/internal/config"
	"github.com/ju4n97/esquema/internal/manifest"
)

// TestGenerateOpenAPI verifies OpenAPI 3.1 specification compilation, validation, and route filtering.
func TestGenerateOpenAPI(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		OpenAPI: config.OpenAPIMetadata{
			Title:   "Inventory Service",
			Version: "1.0.0",
		},
		Schemas: map[string]config.Schema{
			"Product": {
				Name: "Product",
				Fields: map[string]config.Field{
					"id": {
						Name:     "id",
						Type:     config.DataTypeInteger,
						Required: true,
					},
					"sku": {
						Name:     "sku",
						Type:     config.DataTypeString,
						Required: true,
					},
				},
			},
		},
		Endpoints: []config.CompiledEndpoint{
			{
				Method:       "GET",
				Path:         "/products/{id}",
				RoutePattern: "GET /products/{id}",
				Summary:      "Fetch product by identifier",
				Tag:          "products",
				Request: config.RequestRules{
					Path: map[string]config.Field{
						"id": {
							Name:     "id",
							Type:     config.DataTypeInteger,
							Required: true,
						},
					},
				},
				Pipeline: []config.Step{
					{
						Type: config.StepTypeRespond,
						Respond: &config.RespondStep{
							Status:    200,
							SchemaRef: "Product",
						},
					},
				},
			},
			{
				Method:       "GET",
				Path:         "/docs",
				RoutePattern: "GET /docs",
				Pipeline: []config.Step{
					{
						Type: config.StepTypeDocs,
						Docs: &config.DocsStep{Renderer: config.DocsRendererScalar},
					},
				},
			},
			{
				Method:       "GET",
				Path:         "/internal/metrics",
				RoutePattern: "GET /internal/metrics",
				Hidden:       true,
				Pipeline: []config.Step{
					{
						Type:    config.StepTypeRespond,
						Respond: &config.RespondStep{Status: 200},
					},
				},
			},
		},
	}

	rawJSON, err := manifest.GenerateOpenAPI(cfg, string(config.SpecFormatJSON))
	if err != nil {
		t.Fatalf("GenerateOpenAPI failed: %v", err)
	}

	var doc map[string]any
	err = json.Unmarshal(rawJSON, &doc)
	if err != nil {
		t.Fatalf("failed to decode emitted OpenAPI specification: %v", err)
	}

	if doc["openapi"] != "3.1.0" {
		t.Errorf("openapi version = %v; want '3.1.0'", doc["openapi"])
	}

	paths, ok := doc["paths"].(map[string]any)
	if !ok {
		t.Fatalf("paths is not an object: %T", doc["paths"])
	}

	if _, exists := paths["/products/{id}"]; !exists {
		t.Error("expected '/products/{id}' in paths catalog")
	}
	if _, exists := paths["/docs"]; exists {
		t.Error("expected documentation route '/docs' to be excluded from paths")
	}
	if _, exists := paths["/internal/metrics"]; exists {
		t.Error("expected hidden route '/internal/metrics' to be excluded from paths")
	}

	components, ok := doc["components"].(map[string]any)
	if !ok {
		t.Fatalf("components is not an object: %T", doc["components"])
	}

	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		t.Fatalf("schemas is not an object: %T", components["schemas"])
	}

	if _, exists := schemas["Product"]; !exists {
		t.Error("expected 'Product' schema under components.schemas")
	}
	if _, exists := schemas["ProblemDetails"]; !exists {
		t.Error("expected standard 'ProblemDetails' schema under components.schemas")
	}
}
