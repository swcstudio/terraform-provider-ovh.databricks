package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ovh/go-ovh/ovh"
)

var _ provider.Provider = &DatabricksOVHProvider{}

type DatabricksOVHProvider struct {
	version string
}

type DatabricksOVHProviderModel struct {
	OVHEndpoint           types.String `tfsdk:"ovh_endpoint"`
	OVHApplicationKey     types.String `tfsdk:"ovh_application_key"`
	OVHApplicationSecret  types.String `tfsdk:"ovh_application_secret"`
	OVHConsumerKey        types.String `tfsdk:"ovh_consumer_key"`
	OVHProjectID          types.String `tfsdk:"ovh_project_id"`
	DatabricksAccountID   types.String `tfsdk:"databricks_account_id"`
	DatabricksUsername    types.String `tfsdk:"databricks_username"`
	DatabricksPassword    types.String `tfsdk:"databricks_password"`
	DatabricksToken       types.String `tfsdk:"databricks_token"`
}

type Config struct {
	OVHClient *ovh.Client
}

func (p *DatabricksOVHProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "databricks-ovh"
	resp.Version = p.version
}

func (p *DatabricksOVHProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Databricks OVH provider enables management of Databricks resources on OVH cloud infrastructure.",
		Attributes: map[string]schema.Attribute{
			"ovh_endpoint": schema.StringAttribute{
				Description: "OVH API endpoint",
				Required:    true,
			},
			"ovh_application_key": schema.StringAttribute{
				Description: "OVH application key",
				Required:    true,
			},
			"ovh_application_secret": schema.StringAttribute{
				Description: "OVH application secret",
				Required:    true,
				Sensitive:   true,
			},
			"ovh_consumer_key": schema.StringAttribute{
				Description: "OVH consumer key",
				Required:    true,
				Sensitive:   true,
			},
			"ovh_project_id": schema.StringAttribute{
				Description: "OVH project ID",
				Optional:    true,
			},
			"databricks_account_id": schema.StringAttribute{
				Description: "Databricks account ID",
				Optional:    true,
			},
			"databricks_username": schema.StringAttribute{
				Description: "Databricks username",
				Optional:    true,
			},
			"databricks_password": schema.StringAttribute{
				Description: "Databricks password",
				Optional:    true,
				Sensitive:   true,
			},
			"databricks_token": schema.StringAttribute{
				Description: "Databricks personal access token",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *DatabricksOVHProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Databricks OVH client")

	var config DatabricksOVHProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.OVHEndpoint.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown OVH API Endpoint",
			"The provider cannot create the OVH API client as there is an unknown configuration value for the OVH API endpoint.",
		)
	}

	if config.OVHApplicationKey.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown OVH Application Key",
			"The provider cannot create the OVH API client as there is an unknown configuration value for the OVH application key.",
		)
	}

	if config.OVHApplicationSecret.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown OVH Application Secret",
			"The provider cannot create the OVH API client as there is an unknown configuration value for the OVH application secret.",
		)
	}

	if config.OVHConsumerKey.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown OVH Consumer Key",
			"The provider cannot create the OVH API client as there is an unknown configuration value for the OVH consumer key.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := config.OVHEndpoint.ValueString()
	if endpoint == "" {
		endpoint = os.Getenv("OVH_ENDPOINT")
	}

	applicationKey := config.OVHApplicationKey.ValueString()
	if applicationKey == "" {
		applicationKey = os.Getenv("OVH_APPLICATION_KEY")
	}

	applicationSecret := config.OVHApplicationSecret.ValueString()
	if applicationSecret == "" {
		applicationSecret = os.Getenv("OVH_APPLICATION_SECRET")
	}

	consumerKey := config.OVHConsumerKey.ValueString()
	if consumerKey == "" {
		consumerKey = os.Getenv("OVH_CONSUMER_KEY")
	}

	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Missing OVH API Endpoint",
			"The provider requires an OVH API endpoint to be configured.",
		)
	}

	if applicationKey == "" {
		resp.Diagnostics.AddError(
			"Missing OVH Application Key",
			"The provider requires an OVH application key to be configured.",
		)
	}

	if applicationSecret == "" {
		resp.Diagnostics.AddError(
			"Missing OVH Application Secret",
			"The provider requires an OVH application secret to be configured.",
		)
	}

	if consumerKey == "" {
		resp.Diagnostics.AddError(
			"Missing OVH Consumer Key",
			"The provider requires an OVH consumer key to be configured.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ovh_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "ovh_application_key", applicationKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ovh_application_secret")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ovh_consumer_key")

	tflog.Debug(ctx, "Creating OVH client")

	client, err := ovh.NewClient(
		endpoint,
		applicationKey,
		applicationSecret,
		consumerKey,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OVH API Client",
			"An unexpected error occurred when creating the OVH API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OVH Client Error: "+err.Error(),
		)
		return
	}

	providerConfig := &Config{
		OVHClient: client,
	}

	resp.DataSourceData = providerConfig
	resp.ResourceData = providerConfig

	tflog.Info(ctx, "Configured Databricks OVH client", map[string]any{"success": true})
}

func (p *DatabricksOVHProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDatabricksWorkspaceResource,
		NewDatabricksJobResource,
		NewDatabricksNotebookResource,
		NewDatabricksSecretScopeResource,
		NewDatabricksInstancePoolResource,
		NewDatabricksClusterPolicyResource,
	}
}

func (p *DatabricksOVHProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDatabricksWorkspacesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DatabricksOVHProvider{
			version: version,
		}
	}
}
