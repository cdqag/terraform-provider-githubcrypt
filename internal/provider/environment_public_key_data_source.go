package provider

import (
	"context"
	"fmt"

	github "github.com/google/go-github/v64/github"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &EnvironmentPublicKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentPublicKeyDataSource{}
)

func NewEnvironmentPublicKeyDataSource() datasource.DataSource {
	return &EnvironmentPublicKeyDataSource{}
}

// EnvironmentPublicKeyDataSource defines the data source implementation.
type EnvironmentPublicKeyDataSource struct {
	client *github.Client
}

func (d *EnvironmentPublicKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*github.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *github.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

type EnvironmentPublicKeyDataSourceModel struct {
	RepoID      types.Int64  `tfsdk:"repo_id"`
	PublicKey   types.String `tfsdk:"public_key"`
	Environment types.String `tfsdk:"environment"`
}

func (d *EnvironmentPublicKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_public_key"
}

func (d *EnvironmentPublicKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "GitHub Repository Environment Public Key Data Source",

		Attributes: map[string]schema.Attribute{
			"repo_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the GitHub repository.",
				Required:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "The name of the GitHub repository environment.",
				Required:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "The public key of the GitHub repository environment.",
				Required:            true,
				Computed:            true,
			},
		},
	}
}

func (d *EnvironmentPublicKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvironmentPublicKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	key, _, err := d.client.Actions.GetEnvPublicKey(ctx, int(state.RepoID.ValueInt64()), state.Environment.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read repository (%d) public key", state.RepoID.ValueInt64()),
			fmt.Sprintf("An unexpected error occurred while reading the repository (%d) public key...", state.RepoID.ValueInt64())+
				"GitHub Actions API client error: "+err.Error(),
		)
		return
	}

	state = EnvironmentPublicKeyDataSourceModel{
		RepoID:      types.Int64Value(state.RepoID.ValueInt64()),
		PublicKey:   types.StringValue(key.GetKey()),
		Environment: types.StringValue(state.Environment.ValueString()),
	}

	// Save data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
