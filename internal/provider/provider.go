package provider

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v64/github"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure GitHubCryptProvider satisfies various provider interfaces.
var _ provider.Provider = &GitHubCryptProvider{}
var _ provider.ProviderWithFunctions = &GitHubCryptProvider{}

// GitHubCryptProvider defines the provider implementation.
type GitHubCryptProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// GitHubCryptProviderModel describes the provider data model.
type GitHubCryptProviderModel struct {
	Owner             types.String `tfsdk:"owner"`
	AppID             types.String `tfsdk:"app_id"`
	AppInstallationID types.String `tfsdk:"app_installation_id"`
	PemFile           types.String `tfsdk:"pem_file"`
}

func (p *GitHubCryptProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "githubcrypt"
	resp.Version = p.version
}

func (p *GitHubCryptProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				MarkdownDescription: "Owner of the GitHub organization.",
				Optional:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The GitHub App ID.",
				Optional:            true,
				Sensitive:           true,
			},
			"app_installation_id": schema.StringAttribute{
				MarkdownDescription: "The GitHub App Installation ID.",
				Optional:            true,
				Sensitive:           true,
			},
			"pem_file": schema.StringAttribute{
				MarkdownDescription: "GitHub App PEM file.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *GitHubCryptProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config GitHubCryptProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Owner.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("owner"),
			"Unknown GitHub organization owner",
			"The provider cannot create the GitHub API client without a GitHub organization owner",
		)
	}

	if config.AppID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_id"),
			"Unknown GitHub App ID",
			"The provider cannot create the GitHub API client without a GitHub App ID",
		)
	}

	if config.AppInstallationID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_installation_id"),
			"Unknown GitHub App Installation ID",
			"The provider cannot create the GitHub API client without a GitHub App Installation ID",
		)
	}

	if config.PemFile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("pem_file"),
			"Unknown GitHub App PEM file",
			"The provider cannot create the GitHub API client without a GitHub App PEM file",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	owner := os.Getenv("GITHUB_OWNER")
	pemFile := os.Getenv("GITHUB_PEM_FILE")

	appIDString := os.Getenv("GITHUB_APP_ID")
	appID, err := strconv.ParseInt(appIDString, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse GitHub App ID",
			"Error while parsing GitHub App ID from environment variable: "+err.Error(),
		)
	}

	appInstallationIDString := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	appInstallationID, err := strconv.ParseInt(appInstallationIDString, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse GitHub App Installation ID",
			"Error while parsing GitHub App Installation ID from environment variable: "+err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Owner.IsNull() {
		owner = config.Owner.ValueString()
	}

	if !config.AppID.IsNull() {
		appIDInt, err := strconv.ParseInt((config.AppID.ValueString()), 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("app_id"),
				"Failed to parse GitHub App ID",
				"Error while converting the value from config (string) to int64.",
			)
		}
		appID = appIDInt
	}

	if !config.AppInstallationID.IsNull() {
		appInstallationIDInt, err := strconv.ParseInt((config.AppInstallationID.ValueString()), 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("app_installation_id"),
				"Failed to parse GitHub App Installation ID",
				"Error while converting the value from config (string) to int64.",
			)
		}
		appInstallationID = appInstallationIDInt
	}

	if !config.PemFile.IsNull() {
		pemFile = config.PemFile.ValueString()
	}

	// Check if the values are set
	if owner == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("owner"),
			"Missing GitHub organization owner",
			"The provider cannot create the GitHub client without the owner"+
				"Set the `owner` attribute in the provider configuration or set the `GITHUB_OWNER` environment variable.",
		)
	}

	if appID == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_id"),
			"Missing GitHub App ID",
			"The provider cannot create the GitHub client without the GitHub App ID"+
				"Set the `app_id` attribute in the provider configuration or set the `GITHUB_APP_ID` environment variable.",
		)
	}

	if appInstallationID == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_installation_id"),
			"Missing GitHub App Installation ID",
			"The provider cannot create the GitHub client without the GitHub App Installation ID"+
				"Set the `app_installation_id` attribute in the provider configuration or set the `GITHUB_APP_INSTALLATION_ID` environment variable.",
		)
	}

	if pemFile == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("pem_file"),
			"Missing GitHub App PEM file",
			"The provider cannot create the GitHub client without the GitHub App PEM file"+
				"Set the `pem_file` attribute in the provider configuration or set the `GITHUB_PEM_FILE` environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Configure the GitHub API client
	itr, err := ghinstallation.New(http.DefaultTransport, appID, appInstallationID, []byte(pemFile))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create GitHub App installation client",
			"Error while generating RoundTriper transport method for GitHub App authentication"+
				"GitHub App installation client error: "+err.Error(),
		)
		return
	}

	githubClient := github.NewClient(&http.Client{Transport: itr})

	resp.DataSourceData = githubClient
	resp.ResourceData = githubClient
}

func (p *GitHubCryptProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *GitHubCryptProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEnvironmentPublicKeyDataSource,
		NewEncryptedEnvironmentSecretDataSource,
	}
}

func (p *GitHubCryptProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GitHubCryptProvider{
			version: version,
		}
	}
}
