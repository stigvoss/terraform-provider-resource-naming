// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &ResourceNamingProvider{}
var _ provider.ProviderWithFunctions = &ResourceNamingProvider{}

type ResourceNamingProvider struct {
	version         string
	namingFormat    string
	globalVariables map[string]string
}

// ScaffoldingProviderModel describes the provider data model.
type ResourceNamingProviderModel struct {
	NamingFormat    types.String `tfsdk:"naming_format"`
	GlobalVariables types.Map    `tfsdk:"global_variables"`
}

func (p *ResourceNamingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "resource_naming"
	resp.Version = p.version
}

func (p *ResourceNamingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"naming_format": schema.StringAttribute{
				MarkdownDescription: "The format string for resource names. Use `{placeholder}` syntax for dynamic parts, e.g. `\"{resource_type}-{resource_name}-{env}-{region}-{discriminator}\"`.",
				Required:            true,
			},
			"global_variables": schema.MapAttribute{
				MarkdownDescription: "A map of global name variables shared across all resources. These are combined with per-resource dynamic components when generating names.",
				ElementType:         types.StringType,
				Required:            false,
			},
		},
	}
}

func (p *ResourceNamingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ResourceNamingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p.namingFormat = data.NamingFormat.ValueString()
}

func (p *ResourceNamingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *ResourceNamingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *ResourceNamingProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
	//return []func() function.Function{
	//	NewExampleFunction,
	//}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ResourceNamingProvider{
			version: version,
		}
	}
}
