// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"resource_naming": providerserver.NewProtocol6WithError(New("test")()),
}

// providerConfig returns a minimal provider configuration block for use in acceptance tests.
func providerConfig(namingFormat string) string {
	return fmt.Sprintf(`
provider "resource_naming" {
  naming_format = %q
}
`, namingFormat)
}

// ---------------------------------------------------------------------------
// Unit tests — no Terraform CLI required
// ---------------------------------------------------------------------------

func TestProvider_Metadata(t *testing.T) {
	p := New("1.0.0")()

	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), req, resp)

	if resp.TypeName != "resource_naming" {
		t.Errorf("expected TypeName %q, got %q", "resource_naming", resp.TypeName)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected Version %q, got %q", "1.0.0", resp.Version)
	}
}

func TestProvider_Schema_HasRequiredAttributes(t *testing.T) {
	p := New("test")()

	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %s", resp.Diagnostics)
	}

	attrs := resp.Schema.Attributes
	if _, ok := attrs["naming_format"]; !ok {
		t.Error("schema is missing required attribute 'naming_format'")
	}
	if _, ok := attrs["global_variables"]; !ok {
		t.Error("schema is missing attribute 'global_variables'")
	}
}

// ---------------------------------------------------------------------------
// Acceptance tests — require TF_ACC=1
// ---------------------------------------------------------------------------

// TestAccProvider_WithNamingFormat verifies that the provider can be configured
// with only the required naming_format attribute.
func TestAccProvider_WithNamingFormat(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig("{resource_type}-{env}-{region}"),
			},
		},
	})
}

// TestAccProvider_WithGlobalVariables verifies that the provider can be configured
// with both naming_format and global_variables.
func TestAccProvider_WithGlobalVariables(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "resource_naming" {
  naming_format = "{env}-{region}-{app}"
  global_variables = {
    env    = "prod"
    region = "us-east-1"
  }
}
`,
			},
		},
	})
}

// TestAccProvider_MissingNamingFormat verifies that omitting the required
// naming_format attribute produces a configuration error.
func TestAccProvider_MissingNamingFormat(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "resource_naming" {}
`,
				ExpectError: regexp.MustCompile(`The argument "naming_format" is required`),
			},
		},
	})
}

