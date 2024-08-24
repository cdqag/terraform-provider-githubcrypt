package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	github "github.com/google/go-github/v64/github"
	"golang.org/x/crypto/nacl/box"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &EncryptedEnvironmentSecretDataSource{}
	_ datasource.DataSourceWithConfigure = &EncryptedEnvironmentSecretDataSource{}
)

func NewEncryptedEnvironmentSecretDataSource() datasource.DataSource {
	return &EncryptedEnvironmentSecretDataSource{}
}

// EncryptedEnvironmentSecretDataSource defines the data source implementation.
type EncryptedEnvironmentSecretDataSource struct {
	client *github.Client
}

func (d *EncryptedEnvironmentSecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type EncryptedEnvironmentSecretDataSourceModel struct {
	PublicKeyBase64       types.String `tfsdk:"public_key_base64"`
	Secret                types.String `tfsdk:"secret"`
	SecretEncryptedBase64 types.String `tfsdk:"secret_encrypted_base64"`
}

func (d *EncryptedEnvironmentSecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encrypted_environment_secret"
}

func (d *EncryptedEnvironmentSecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source encrypts a secret with the public key of a GitHub repository environment.",

		Attributes: map[string]schema.Attribute{
			"public_key_base64": schema.StringAttribute{
				MarkdownDescription: "The public key of the GitHub repository environment encoded in Base64",
				Required:            true,
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "The secret to encrypt.",
				Required:            true,
			},
			"secret_encrypted_base64": schema.StringAttribute{
				MarkdownDescription: "The encrypted secret in base64.",
				Computed:            true,
			},
		},
	}
}

func (d *EncryptedEnvironmentSecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EncryptedEnvironmentSecretDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Decode publick key
	decodedKey, err := base64.StdEncoding.DecodeString(state.PublicKeyBase64.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to crypt the secret with repository environment public key."),
			fmt.Sprintf("Failed to decode repository's environment's public key")+
				"Base64 error: "+err.Error(),
		)
		return
	}

	// Encrypt secret
	var keyBytes [32]byte
	copy(keyBytes[:], decodedKey)

	var encryptedBytes []byte
	encSec, err := box.SealAnonymous(encryptedBytes, []byte(state.Secret.ValueString()), &keyBytes, nil)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to crypt the secret with repository environment public key."),
			fmt.Sprintf("Failed to encrypt the secret with the public key.")+
				"Error: "+err.Error(),
		)
		return
	}

	// Encode encrypted secret
	encSecEncoded := base64.StdEncoding.EncodeToString(encSec)

	// Set values into state
	state = EncryptedEnvironmentSecretDataSourceModel{
		PublicKeyBase64:       types.StringValue(state.PublicKeyBase64.ValueString()),
		Secret:                types.StringValue(state.Secret.ValueString()),
		SecretEncryptedBase64: types.StringValue(encSecEncoded),
	}

	// Save data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
