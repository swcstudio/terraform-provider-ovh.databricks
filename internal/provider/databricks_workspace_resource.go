package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DatabricksWorkspaceResource{}
var _ resource.ResourceWithImportState = &DatabricksWorkspaceResource{}

func NewDatabricksWorkspaceResource() resource.Resource {
	return &DatabricksWorkspaceResource{}
}

type DatabricksWorkspaceResource struct {
	client *Config
}

type DatabricksWorkspaceResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Region                   types.String `tfsdk:"region"`
	Tier                     types.String `tfsdk:"tier"`
	DeploymentName           types.String `tfsdk:"deployment_name"`
	AWSRegion                types.String `tfsdk:"aws_region"`
	CredentialsID            types.String `tfsdk:"credentials_id"`
	StorageConfigurationID   types.String `tfsdk:"storage_configuration_id"`
	NetworkID                types.String `tfsdk:"network_id"`
	CustomerManagedKeyID     types.String `tfsdk:"customer_managed_key_id"`
	PricingTier              types.String `tfsdk:"pricing_tier"`
	CustomTags               types.Map    `tfsdk:"custom_tags"`
	OVHOptimization          types.Bool   `tfsdk:"ovh_optimization"`
	CostTracking             types.Bool   `tfsdk:"cost_tracking"`
	WorkspaceID              types.String `tfsdk:"workspace_id"`
	WorkspaceURL             types.String `tfsdk:"workspace_url"`
	WorkspaceStatus          types.String `tfsdk:"workspace_status"`
	CreationTime             types.String `tfsdk:"creation_time"`
}

func (r *DatabricksWorkspaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *DatabricksWorkspaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Databricks workspace on OVH infrastructure",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Workspace identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Workspace name",
				Required:    true,
			},
			"region": schema.StringAttribute{
				Description: "OVH region",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tier": schema.StringAttribute{
				Description: "Databricks tier",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
			},
			"deployment_name": schema.StringAttribute{
				Description: "Deployment name",
				Optional:    true,
				Computed:    true,
			},
			"aws_region": schema.StringAttribute{
				Description: "AWS region for workspace",
				Optional:    true,
				Computed:    true,
			},
			"credentials_id": schema.StringAttribute{
				Description: "Credentials ID",
				Optional:    true,
			},
			"storage_configuration_id": schema.StringAttribute{
				Description: "Storage configuration ID",
				Optional:    true,
			},
			"network_id": schema.StringAttribute{
				Description: "Network configuration ID",
				Optional:    true,
			},
			"customer_managed_key_id": schema.StringAttribute{
				Description: "Customer managed key ID",
				Optional:    true,
			},
			"pricing_tier": schema.StringAttribute{
				Description: "Pricing tier",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
			},
			"custom_tags": schema.MapAttribute{
				Description: "Custom tags",
				Optional:    true,
				ElementType: types.StringType,
			},
			"ovh_optimization": schema.BoolAttribute{
				Description: "Enable OVH infrastructure optimization",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"cost_tracking": schema.BoolAttribute{
				Description: "Enable cost tracking",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"workspace_id": schema.StringAttribute{
				Description: "Workspace ID",
				Computed:    true,
			},
			"workspace_url": schema.StringAttribute{
				Description: "Workspace URL",
				Computed:    true,
			},
			"workspace_status": schema.StringAttribute{
				Description: "Workspace status",
				Computed:    true,
			},
			"creation_time": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
		},
	}
}

func (r *DatabricksWorkspaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *DatabricksWorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabricksWorkspaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating databricks workspace resource")

	workspaceConfig := map[string]interface{}{
		"name":                     data.Name.ValueString(),
		"region":                   data.Region.ValueString(),
		"tier":                     data.Tier.ValueString(),
		"deploymentName":           data.DeploymentName.ValueString(),
		"awsRegion":                data.AWSRegion.ValueString(),
		"credentialsId":            data.CredentialsID.ValueString(),
		"storageConfigurationId":   data.StorageConfigurationID.ValueString(),
		"networkId":                data.NetworkID.ValueString(),
		"customerManagedKeyId":     data.CustomerManagedKeyID.ValueString(),
		"pricingTier":              data.PricingTier.ValueString(),
		"ovhOptimization":          data.OVHOptimization.ValueBool(),
		"costTracking":             data.CostTracking.ValueBool(),
	}

	if !data.CustomTags.IsNull() {
		var tags map[string]string
		resp.Diagnostics.Append(data.CustomTags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		workspaceConfig["customTags"] = tags
	}

	var result map[string]interface{}
	err := r.client.OVHClient.Post("/cloud/project/databricks/workspace", workspaceConfig, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace, got error: %s", err))
		return
	}

	workspaceId := result["id"].(string)
	data.ID = types.StringValue(workspaceId)

	tflog.Trace(ctx, "created databricks workspace resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksWorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabricksWorkspaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var workspace map[string]interface{}
	err := r.client.OVHClient.Get(fmt.Sprintf("/cloud/project/databricks/workspace/%s", data.ID.ValueString()), &workspace)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace, got error: %s", err))
		return
	}

	if name, ok := workspace["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if region, ok := workspace["region"].(string); ok {
		data.Region = types.StringValue(region)
	}
	if tier, ok := workspace["tier"].(string); ok {
		data.Tier = types.StringValue(tier)
	}
	if deploymentName, ok := workspace["deploymentName"].(string); ok {
		data.DeploymentName = types.StringValue(deploymentName)
	}
	if awsRegion, ok := workspace["awsRegion"].(string); ok {
		data.AWSRegion = types.StringValue(awsRegion)
	}
	if credentialsId, ok := workspace["credentialsId"].(string); ok {
		data.CredentialsID = types.StringValue(credentialsId)
	}
	if storageConfigurationId, ok := workspace["storageConfigurationId"].(string); ok {
		data.StorageConfigurationID = types.StringValue(storageConfigurationId)
	}
	if networkId, ok := workspace["networkId"].(string); ok {
		data.NetworkID = types.StringValue(networkId)
	}
	if customerManagedKeyId, ok := workspace["customerManagedKeyId"].(string); ok {
		data.CustomerManagedKeyID = types.StringValue(customerManagedKeyId)
	}
	if pricingTier, ok := workspace["pricingTier"].(string); ok {
		data.PricingTier = types.StringValue(pricingTier)
	}
	if ovhOptimization, ok := workspace["ovhOptimization"].(bool); ok {
		data.OVHOptimization = types.BoolValue(ovhOptimization)
	}
	if costTracking, ok := workspace["costTracking"].(bool); ok {
		data.CostTracking = types.BoolValue(costTracking)
	}
	if workspaceId, ok := workspace["workspaceId"].(string); ok {
		data.WorkspaceID = types.StringValue(workspaceId)
	}
	if workspaceUrl, ok := workspace["workspaceUrl"].(string); ok {
		data.WorkspaceURL = types.StringValue(workspaceUrl)
	}
	if workspaceStatus, ok := workspace["workspaceStatus"].(string); ok {
		data.WorkspaceStatus = types.StringValue(workspaceStatus)
	}
	if creationTime, ok := workspace["creationTime"].(string); ok {
		data.CreationTime = types.StringValue(creationTime)
	}

	if tags, ok := workspace["customTags"].(map[string]interface{}); ok {
		tagMap := make(map[string]string)
		for k, v := range tags {
			if str, ok := v.(string); ok {
				tagMap[k] = str
			}
		}
		tagsValue, diags := types.MapValueFrom(ctx, types.StringType, tagMap)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			data.CustomTags = tagsValue
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksWorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabricksWorkspaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateConfig := map[string]interface{}{}

	if !data.Name.IsNull() {
		updateConfig["name"] = data.Name.ValueString()
	}
	if !data.Tier.IsNull() {
		updateConfig["tier"] = data.Tier.ValueString()
	}
	if !data.PricingTier.IsNull() {
		updateConfig["pricingTier"] = data.PricingTier.ValueString()
	}
	if !data.CustomTags.IsNull() {
		var tags map[string]string
		resp.Diagnostics.Append(data.CustomTags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateConfig["customTags"] = tags
	}

	err := r.client.OVHClient.Put(fmt.Sprintf("/cloud/project/databricks/workspace/%s", data.ID.ValueString()), updateConfig, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workspace, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabricksWorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabricksWorkspaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OVHClient.Delete(fmt.Sprintf("/cloud/project/databricks/workspace/%s", data.ID.ValueString()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workspace, got error: %s", err))
		return
	}
}

func (r *DatabricksWorkspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
